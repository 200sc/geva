package population

import (
	"goevo/neural"
	"math"
	"sort"
)

type Population struct {
	Members      []neural.Network
	Mutation     *neural.NetworkGenerationOptions
	Size         int
	Selection    SelectionMethod
	Crossover    CrossoverMethod
	Pairing      PairingMethod
	TestInputs   [][]float64
	TestExpected [][]float64
	Elites       int
	Fitnesses    []int
	LowFitness   int
	MaxFitness   int
}

// This will change as more things take place
// in a generation. Selection, Crossover, and Mutation
// are granted.
func (p_p *Population) NextGeneration() *Population {
	p := *p_p
	// The number of parents in the next generation
	parentSize := p.Size / p.Selection.GetParentProportion()

	p = *(p.GenerateFitness())
	elites := p.GetElites()
	nextGen := p.Selection.Select(&p)

	// Ensure that the elites (the best members)
	// stay in the next generation
	for i, elite := range elites {
		nextGen[i+parentSize] = nextGen[i]
		nextGen[i] = elite
	}
	parentSize += p.Elites

	// Determine our pairing method and then
	// crossover pairs for children in the next generation.
	pairs := p.Pairing.Pair(nextGen, parentSize)
	p.Members = p.Crossover.Crossover(nextGen, parentSize, pairs)

	// The elites are not subject to mutation.
	for i := p.Elites; i < len(p.Members); i++ {
		p.Members[i] = *p.Mutation.Mutate(p.Members[i])
	}
	return &p
}

func (p_p *Population) GenerateFitness() *Population {
	p := *p_p

	channels := make([]chan int, p.Size)

	for i := 0; i < p.Size; i++ {
		channels[i] = make(chan int)

		go func(n *neural.Network, ch chan int, inputs [][]float64, expected [][]float64) {
			ch <- (*n).Fitness(inputs, expected)
		}(&(p.Members[i]), channels[i], p.TestInputs, p.TestExpected)
	}

	p.LowFitness = math.MaxInt32
	p.MaxFitness = 0

	for i := 0; i < p.Size; i++ {
		v := <-channels[i]
		close(channels[i])
		if v < p.LowFitness {
			p.LowFitness = v
		} else if v > p.MaxFitness {
			p.MaxFitness = v
		}
		p.Fitnesses[i] = v
	}

	return &p
}

func (p_p *Population) GetElites() []neural.Network {
	p := *p_p

	fitMap := make(map[int][]int)
	elites := make([]neural.Network, p.Elites)

	for i := 0; i < p.Size; i++ {
		f := p.Fitnesses[i]
		if v, ok := fitMap[f]; ok {
			fitMap[f] = append(v, i)
		} else {
			fitMap[f] = []int{i}
		}
	}

	keys := KeySet_Int_SlInt(fitMap)
	sort.Ints(keys)
	i := 0
	j := 0
	for i < p.Elites {
		for k := 0; k < len(fitMap[keys[j]]); k++ {
			if i >= p.Elites {
				return elites
			}
			elites[i] = p.Members[fitMap[keys[j]][k]]
			i++
		}
		j++
	}
	return elites
}

func KeySet_Int_SlInt(m map[int][]int) []int {
	keys := make([]int, len(m))

	i := 0
	for k := range m {
		keys[i] = k
		i++
	}

	return keys
}

func (p_p *Population) Weights(power float64) ([]float64, []float64) {
	p := *p_p
	fitnesses := p.Fitnesses
	maxFitness := p.MaxFitness

	weights := make([]float64, len(fitnesses))
	cumulativeWeights := make([]float64, len(fitnesses))

	// Transform values which are low to equivalent high
	// values on the same scale, applying the power
	// as a further bias scaling towards the best
	// individuals.
	for i := 0; i < len(fitnesses); i++ {
		weights[i] = math.Pow(float64((fitnesses[i]*-1)+maxFitness+1), power)
	}

	cumulativeWeights[0] = weights[0]

	for i := 0; i < len(fitnesses)-1; i++ {
		cumulativeWeights[i+1] = cumulativeWeights[i] + weights[i+1]
	}

	return weights, cumulativeWeights
}

func (p_p *Population) Print() {
	for _, v := range p_p.Members {
		v.Print()
	}
}

// Used as Generic-esque helpers for populations

type SelectionMethod interface {
	Select(p_p *Population) []neural.Network
	GetParentProportion() int
}

type CrossoverMethod interface {
	Crossover(nn []neural.Network, populated int, pairs [][]int) []neural.Network
}

type PairingMethod interface {
	Pair(nn []neural.Network, populated int) [][]int
}
