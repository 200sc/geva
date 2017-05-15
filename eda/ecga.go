package eda

import (
	"container/heap"
	"math"
	"math/rand"

	"bitbucket.org/StephenPatrick/goevo/eda/stat"
	"bitbucket.org/StephenPatrick/goevo/env"
	"bitbucket.org/StephenPatrick/goevo/pop"
)

// ECGA represents the Extended Compact Genetic Algorithm
type ECGA struct {
	Base
	P      *pop.Population
	Blocks [][]int
}

// Adjust for an ecga creates an ecga population based on the ecga's
// understanding of its building blocks and then refreshes those
// building blocks
func (ecga *ECGA) Adjust() Model {
	// Generate pop from mpm model
	ecga.P = ecga.ECGAPop()
	ecga.UpdateMDM()
	//fmt.Println(ecga.Blocks)
	return ecga
}

// Mutate on an ecga mutates each element of the ecga population
func (ecga *ECGA) Mutate() {
	for _, m := range ecga.P.Members {
		m.(*EnvInd).F.Mutate(ecga.mutationRate, ecga.fmutator)
	}
}

// UpdateMDM updates the building blocks and environment of an ecga
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

// ECGAModel returns an initialized ECGA EDA
func ECGAModel(opts ...Option) (Model, error) {
	var err error
	ecga := new(ECGA)
	ecga.Base, err = DefaultBase(opts...)
	// Random Pop
	ecga.P = ecga.Pop()
	ecga.UpdateMDM()
	return ecga, err
}

// ECGAPop returns a population where some members are from ecga.P
// and some are sampled from ecga.Blocks and ecga.P at random. The
// proportion of new to old members is based on learningRate.
func (ecga *ECGA) ECGAPop() *pop.Population {
	newMemberCt := int(float64(len(ecga.P.Members)) * ecga.learningRate)
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
	ecga.ReplaceLowFitnesses(ecga.P, newMembers)
	return ecga.P
}

// MDMModel refreshes the ecga's building blocks
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
	var outBlocks [][]int
	for i, b := range blocks {
		if _, ok := merged[i]; !ok {
			outBlocks = append(outBlocks, b)
		}
	}
	ecga.Blocks = outBlocks
}

// BlockComplexity returns the complexity of a given block definition
func (ecga *ECGA) BlockComplexity(envs []*env.F, b []int) float64 {
	return (math.Log2(float64(ecga.samples)) * ecga.ModelComplexity(b)) +
		(float64(ecga.samples) * ecga.CombinedComplexity(envs, b))
}

// ModelComplexity punishes blocks exponentially for being long
// in effect, this means ecgas can't develop building blocks longer than
// maybe four or five indices.
func (ecga *ECGA) ModelComplexity(b []int) float64 {
	return math.Pow(2, float64(len(b)))
}

// CombinedComplexity returns the sum marginal complexity of the
// block's indices in the sample set
func (ecga *ECGA) CombinedComplexity(envs []*env.F, b []int) float64 {
	return stat.MarginalEntropy(envs, b)
}

// Papers : https://pdfs.semanticscholar.org/eeee/a9fdade929cb3fc9a99631d3541ef7005079.pdf
// http://www.kumarasastry.com/wp-content/files/2000026.pdf
