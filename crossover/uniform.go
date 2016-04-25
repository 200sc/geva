package crossover

import (
	"goevo/neural"
	"math/rand"
)

// Choose a bunch of random neurons from each network
// and make a new network out of them.
// I don't think this is a very good idea for neural
// networks, but we'll see.
type UniformCrossover struct {
	// This proportion of neurons that are chosen
	// from the first network selected.
	// The remaining proporiton 1 - chosenProportion
	// come from the other network.
	// Cannot be negative.
	ChosenProportion float64
}

func (pc_p UniformCrossover) Crossover(nn []neural.Network, populated int, pairs [][]int) []neural.Network {

	pairIndex := 0

	for j := populated; j < len(nn); j++ {

		n1 := nn[pairs[pairIndex][0]].Body
		n2 := nn[pairs[pairIndex][1]].Body

		newBody := n1.CopyStructure()
		// This assumes that each network has the same dimensions!
		for i := 0; i < len(n1); i++ {
			for k := 0; k < len(n1[i]); k++ {
				if rand.Float64() < pc_p.ChosenProportion {
					// This assumes our value is copied and is not just a pointer.
					newBody[i][k] = n1[i][k]
				} else {
					newBody[i][k] = n2[i][k]
				}
			}
		}
		nn[j] = neural.Network{
			Body:      newBody,
			Activator: nn[pairs[pairIndex][0]].Activator,
		}

		pairIndex++
	}

	return nn
}
