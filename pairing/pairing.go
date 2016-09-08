package pairing

import (
	"goevo/population"
	"math/rand"
)

type RandomPairing struct{}

func (rp RandomPairing) Pair(inds []population.Individual, populated int) [][]int {

	out := make([][]int, len(inds)-populated)

	for i := 0; i < len(inds)-populated; i++ {

		index1 := rand.Intn(populated)
		index2 := rand.Intn(populated)

		if index1 == index2 {
			index2 = (index2 + 1) % populated
		}

		out[i] = []int{index1, index2}
	}

	return out
}
