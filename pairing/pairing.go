package pairing

import (
	"goevo/neural"
	"math/rand"
)

type RandomPairing struct{}

func (rp RandomPairing) Pair(nn []neural.Network, populated int) [][]int {

	out := make([][]int, len(nn)-populated)

	for i := 0; i < len(nn)-populated; i++ {

		index1 := rand.Intn(populated)
		index2 := rand.Intn(populated)

		if index1 == index2 {
			index2 = (index2 + 1) % populated
		}

		out[i] = []int{index1, index2}
	}

	return out
}
