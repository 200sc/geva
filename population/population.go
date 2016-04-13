package population

import (
	//"fmt"
	"goevo/neural"
)

type Population struct {
	Members           []neural.Network
	Type              NetworkType
	Size              int
	Selection         SelectionMethod
	Crossover         CrossoverMethod
	TestInputs        [][]float64
	TestExpected      [][]float64
	Weights           []int //unused right now
	CumulativeWeights []int // ^
}

func (p_p *Population) NextGeneration() *Population {
	p := *p_p

	nextGen := p.Selection.Select(p_p)
	p.Members = p.Crossover.Crossover(nextGen, p.Size/p.Selection.GetParentProportion())
	for _, v := range p.Members {
		v = p.Type.Mutate(v)
	}
	return &p
}

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
