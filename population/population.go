package population

import (
	//"fmt"
	"goevo/neural"
	"math"
)

type Population struct {
	Members      []neural.Network
	Type         NetworkType
	Size         int
	Selection    SelectionMethod
	Crossover    CrossoverMethod
	TestInputs   [][]float64
	TestExpected [][]float64
}

// This will change as more things take place
// in a generation. Selection, Crossover, and Mutation
// are granted.
func (p_p *Population) NextGeneration() *Population {
	p := *p_p

	nextGen := p.Selection.Select(p_p)
	p.Members = p.Crossover.Crossover(nextGen, p.Size/p.Selection.GetParentProportion())
	for _, v := range p.Members {
		v = p.Type.Mutate(v)
	}
	return &p
}

// Relative misnomer
// This doesn't calculate the fitness of the population,
// at least not immediately.
// It starts a bunch of goroutines which will then eventually get their fitnesses
// back to you via the channels this returns.

func (p_p *Population) Fitness() []chan int {
	p := *p_p

	channels := make([]chan int, p.Size)

	for i := 0; i < p.Size; i++ {
		channels[i] = make(chan int)

		go func(n *neural.Network, ch chan int, inputs [][]float64, expected [][]float64) {
			ch <- (*n).Fitness(inputs, expected)
		}(&(p.Members[i]), channels[i], p.TestInputs, p.TestExpected)
	}
	return channels
}

func (p_p *Population) Weights(power float64) ([]float64, []float64) {
	p := *p_p

	fitnessChannels := p_p.Fitness()
	fitnesses := make([]int, p.Size)
	weights := make([]float64, p.Size)
	cumulativeWeights := make([]float64, p.Size)

	maxFitness := 0

	for i := 0; i < p.Size; i++ {
		v := <-fitnessChannels[i]
		if v > maxFitness {
			maxFitness = v
		}
		fitnesses[i] = v
	}

	// Transform values which are low to equivalent high
	// values on the same scale, applying the power
	// as a further bias scaling towards the best
	// individuals.
	for i := 0; i < p.Size; i++ {
		weights[i] = math.Pow(float64((fitnesses[i]*-1)+maxFitness+1), power)
	}

	cumulativeWeights[0] = weights[0]

	for i := 0; i < p.Size-1; i++ {
		cumulativeWeights[i+1] = cumulativeWeights[i] + weights[i+1]
	}

	return weights, cumulativeWeights
}

func (p_p *Population) Print() {
	for _, v := range p_p.Members {
		v.Print()
	}
}

// func RouletteSlice(sl []int) []int {
// 	// We want the minimum element (the best element)
// 	// to have len(sl) weight in the roulette. Each
// 	// element in the sorted list indexed at i should
// 	// have len(sl) - i weight.
// 	sort.Ints(sl)
// }

// Used as Generic-esque helpers for populations

type SelectionMethod interface {
	Select(p_p *Population) []neural.Network
	GetParentProportion() int
}

type CrossoverMethod interface {
	Crossover(nn []neural.Network, populated int) []neural.Network
}

type NetworkType interface {
	Generate() neural.Network
	Mutate(neural.Network) neural.Network
}
