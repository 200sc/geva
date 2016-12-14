package pairing

import (
	"container/heap"
	"goevo/population"
	"math/rand"
)

type Random struct{}

func (rp Random) Pair(p *population.Population, populated int) [][]int {

	out := make([][]int, len(p.Members)-populated)

	for i := 0; i < len(p.Members)-populated; i++ {

		index1 := rand.Intn(populated)
		index2 := rand.Intn(populated)

		if index1 == index2 {
			index2 = (index2 + 1) % populated
		}

		out[i] = []int{index1, index2}
	}

	return out
}

// Alpha pairing makes sure every new child has
// at least one parent who was one of the best
// of the previous generation.
type Alpha struct {
	AlphaCount int
}

func (ap Alpha) Pair(p *population.Population, populated int) [][]int {

	out := make([][]int, len(p.Members)-populated)

	bestMembers := make([]int, ap.AlphaCount)

	h := &ValIndexHeap{}
	heap.Init(h)
	for i := 0; i < populated; i++ {
		f := p.Fitnesses[i]
		heap.Push(h, [2]int{f, i})
	}
	for i := 0; i < ap.AlphaCount; i++ {
		valIndex := heap.Pop(h).([2]int)
		bestMembers[i] = valIndex[1]
	}

	for i := 0; i < len(p.Members)-populated; i++ {

		index1 := bestMembers[i%ap.AlphaCount]
		index2 := rand.Intn(populated)

		if index1 == index2 {
			index2 = (index2 + 1) % populated
		}

		out[i] = []int{index1, index2}
	}

	return out
}

type Method interface {
	Pair(p *population.Population, populated int) [][]int
}
