package eda

import (
	"container/heap"
	"math"
	"math/rand"

	"bitbucket.org/StephenPatrick/goevo/eda/stat"
	"bitbucket.org/StephenPatrick/goevo/env"
	"bitbucket.org/StephenPatrick/goevo/pop"
)

type ECGA struct {
	Base
	P                *pop.Population
	Blocks           [][]int
	ParentProportion float64
}

func (ecga *ECGA) Adjust() Model {
	// Generate pop from mpm model
	ecga.P = ecga.ECGAPop()
	ecga.UpdateMDM()
	return ecga
}

func (ecga *ECGA) UpdateMDM() {
	// Selection
	selected := ecga.SelectLearning(ecga.P)
	// Update ecga.F to be the average of the selected members
	ecga.F.SetAll(0.0)
	for _, s := range selected {
		ecga.F.AddF(s.(*EnvInd).F)
	}
	ecga.F.Divide(float64(len(selected)))
	// Update the MPM Model using the selected members
	ecga.MDMModel(selected)
}

func ECGAModel(opts ...Option) (Model, error) {
	var err error
	ecga := new(ECGA)
	ecga.ParentProportion = 0.5
	ecga.Base, err = DefaultBase(opts...)
	// Random Pop
	ecga.P = ecga.Pop()
	ecga.UpdateMDM()
	return ecga, err
}

func (ecga *ECGA) ECGAPop() *pop.Population {
	newMemberCt := int(float64(len(ecga.P.Members)) * (1 - ecga.ParentProportion))
	// Sample from the existing population to generate new
	// members, but sample in blocks.
	// Paper 2 suggests that random sampling of the existing
	// population is equivalent to proportionally selecting.
	newMembers := make([]*env.F, newMemberCt)
	for i := range newMembers {
		newMembers[i] = ecga.F.Copy()
		for _, b := range ecga.Blocks {
			sampleFrom := ecga.P.Members[rand.Intn(len(ecga.P.Members))].(*EnvInd)
			for _, j := range b {
				newMembers[i].Set(j, sampleFrom.F.Get(j))
			}
		}
	}
	// replace the low fitness members in the original population

	// Cast members to Environments
	envs := make([]*env.F, len(ecga.P.Members))
	for i, s := range ecga.P.Members {
		envs[i] = s.(*EnvInd).F
	}

	// Sort envs by fitness
	_, envs = SampleFitnesses(ecga, envs)
	i := 0
	for j := len(envs) - 1; j >= 0; j-- {
		envs[j] = newMembers[i]
		i++
		if i >= len(newMembers) {
			break
		}
	}

	// set the population to be envs cast back to members
	for i := range ecga.P.Members {
		ecga.P.Members[i] = &EnvInd{envs[i]}
	}
	return ecga.P
}

func (ecga *ECGA) MDMModel(selected []pop.Individual) {
	// Cast selected to Environments
	envs := make([]*env.F, len(selected))
	for i, s := range selected {
		envs[i] = s.(*EnvInd).F
	}
	// Blocks are initially all independant
	blocks := make([][]int, ecga.length)
	// At the first iteration, the blocks we should try to merge are all blocks
	newBlocks := make([]int, len(blocks))
	for i := range blocks {
		blocks[i] = []int{i}
		newBlocks[i] = i
	}
	// Now we merge blocks so long as the evaluated MPM is better.
	hp := &ComplexHeap{}
	heap.Init(hp)
	// Evaluate all pairs to merge
	for i, b := range blocks {
		iComplexity := ecga.BlockComplexity(envs, b)
		for j := i + 1; j < len(blocks); j++ {
			jComplexity := ecga.BlockComplexity(envs, blocks[j])
			blockIJ := append(b, blocks[j]...)
			ijComplexity := ecga.BlockComplexity(envs, blockIJ)
			diff := ijComplexity - (jComplexity + iComplexity)
			if diff < 0 {
				heap.Push(hp, ComplexDiff{i, j, diff})
				break
			}
		}
	}
	// Merge the best pair found repeatedly
	merged := make(map[int]bool)
	for len(*hp) != 0 {
		toMerge := heap.Pop(hp).(ComplexDiff)
		i, j := toMerge.i, toMerge.j
		blocks[i] = append(blocks[i], blocks[j]...)
		// This map is used and filtered on later to avoid judging the same
		// blocks twice, an alternative would be to immediately remove blocks[j]
		// and then when performing the filter, all indices which used to be above
		// j would be reduced by one.
		merged[j] = true
		// Filter all elements from the complexity difference heap
		// that contain either i or j.
		newHeap := &ComplexHeap{}
		for len(*hp) != 0 {
			v := heap.Pop(hp).(ComplexDiff)
			if v.i == i || v.i == j || v.j == i || v.j == j {
				continue
			}
			heap.Push(newHeap, v)
		}
		hp = newHeap
		// Add new elements to the complexity difference heap for the new
		// merged index
		iComplexity := ecga.BlockComplexity(envs, blocks[i])
		for j := i + 1; j < len(blocks); j++ {
			if _, ok := merged[j]; !ok {
				jComplexity := ecga.BlockComplexity(envs, blocks[j])
				blockIJ := append(blocks[i], blocks[j]...)
				ijComplexity := ecga.BlockComplexity(envs, blockIJ)
				diff := ijComplexity - (jComplexity + iComplexity)
				if diff < 0 {
					heap.Push(hp, ComplexDiff{i, j, diff})
					break
				}
			}
		}
	}
	// Filter out all merged blocks (whose indices were not used by the resulting
	// merger)
	outBlocks := make([][]int, 0)
	for i, b := range blocks {
		if _, ok := merged[i]; !ok {
			outBlocks = append(outBlocks, b)
		}
	}
	ecga.Blocks = outBlocks
}

func (ecga *ECGA) BlockComplexity(envs []*env.F, b []int) float64 {
	return (math.Log2(float64(ecga.samples)) * ecga.BlockModelComplexity(b)) +
		(float64(ecga.samples) * ecga.BlockCombinedComplexity(envs, b))
}

func (ecga *ECGA) EvalMPM(envs []*env.F, blocks [][]int) float64 {
	return ecga.ModelComplexity(blocks) + ecga.CompressedComplexity(envs, blocks)
}

func (ecga *ECGA) ModelComplexity(blocks [][]int) float64 {
	// There's a good risk of this overflowing for large building blocks
	freqTotal := 0.0
	for _, b := range blocks {
		// catch incremental overflow
		inc := ecga.BlockModelComplexity(b)
		if inc < 0 {
			return math.MaxFloat64
		}
		freqTotal += inc
		// catch sum overflow
		if freqTotal < 0 {
			return math.MaxFloat64
		}
	}
	return math.Log2(float64(ecga.samples)) * freqTotal
}

func (ecga *ECGA) BlockModelComplexity(b []int) float64 {
	return math.Pow(2, float64(len(b)))
}

func (ecga *ECGA) BlockCombinedComplexity(envs []*env.F, b []int) float64 {
	return stat.MarginalEntropy(envs, b)
}

func (ecga *ECGA) CompressedComplexity(envs []*env.F, blocks [][]int) float64 {
	eTotal := 0.0
	for _, b := range blocks {
		eTotal += ecga.BlockCombinedComplexity(envs, b)
	}
	return float64(ecga.samples) * eTotal
}

// ECGA isn't terribly hard to do-- we already have crossover and selection
// covered, and elitism is fine but we should probably give the other models
// who use populations elitism if we give ECGA elitism. The one hard thing
// is the MPM model which in itself isn't hard but requires that we split
// the problem space up into building blocks, and requires that we learn these
// building blocks, I guess, although I'm not certain on that
//
// Yes, according to paper 1 we begin the MDM model with the assumption that
// every index is its own building block, and sees if merging all pairs is
// helpful.
//
// Papers : https://pdfs.semanticscholar.org/eeee/a9fdade929cb3fc9a99631d3541ef7005079.pdf
// http://www.kumarasastry.com/wp-content/files/2000026.pdf
