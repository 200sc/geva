package goevo

import (
	"goevo/neural"
)

type Population struct {
	Members           []neural.Network
	generationOptions *neural.NetworkGenerationOptions
	Size              int
	selection         SelectionMethod
	crossover         CrossoverMethod
	testInputs        [][]bool
	testExpected      [][]bool
	weights           []int
	cumulativeweights []int
}

func (p_p *Population) NextGeneration() *Population {
	p := *p_p

	nextGen := p.selection.Select(p_p)
	p.Members = p.crossover.Crossover(nextGen)
	return &p
}

func (p_p *Population) Fitness() []chan int {
	p := *p_p

	channels := make([]chan int, p.Size)

	for i := 0; i < p.Size; i++ {

		go func(n *neural.Network, ch chan int, inputs [][]bool, expected [][]bool) {
			ch <- (*n).Fitness(inputs, expected)

		}(&(p.Members[i]), channels[i], p.testInputs, p.testExpected)
	}

	return channels
}

// func RouletteSlice(sl []int) []int {
// 	// We want the minimum element (the best element)
// 	// to have len(sl) weight in the roulette. Each
// 	// element in the sorted list indexed at i should
// 	// have len(sl) - i weight.
// 	sort.Ints(sl)
// }

type SelectionMethod interface {
	Select(p_p *Population) []neural.Network
}

type CrossoverMethod interface {
	Crossover(nn []neural.Network) []neural.Network
}
