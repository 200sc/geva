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

func (pc_p UniformCrossover) Crossover(nn []neural.Network, populated int) []neural.Network {

	for j := populated; j < len(nn); j++ {

		// In the future, the actual method for selecting
		// pairs to crossover should be variable.
		// Here it is random.
		index1 := rand.Intn(populated)
		index2 := rand.Intn(populated)

		if index1 == index2 {
			index2 = (index2 + 1) % populated
		}

		n1 := nn[index1]
		n2 := nn[index2]

		nn[j] = n1.CopyStructure()
		// This assumes that each network has the same dimensions!
		for i := 0; i < n1.Length(); i++ {
			for k := 0; k < n1.ColLength(i); k++ {
				if rand.Float64() < pc_p.ChosenProportion {
					// This assumes our value is copied.
					nn[j].Set(i, k, n1.Get(i, k))
				} else {
					nn[j].Set(i, k, n2.Get(i, k))
				}
			}
		}
	}

	return nn
}
