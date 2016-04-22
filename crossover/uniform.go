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

func (pc_p UniformCrossover) Crossover(nn []neural.ModularNetwork, populated int) []neural.ModularNetwork {

	for j := populated; j < len(nn); j++ {

		// In the future, the actual method for selecting
		// pairs to crossover should be variable.
		// Here it is random.
		index1 := rand.Intn(populated)
		index2 := rand.Intn(populated)

		if index1 == index2 {
			index2 = (index2 + 1) % populated
		}

		n1 := nn[index1].Body
		n2 := nn[index2].Body

		newBody := n1.CopyStructure()
		// This assumes that each network has the same dimensions!
		for i := 0; i < len(n1); i++ {
			for k := 0; k < len(n1[i]); k++ {
				if rand.Float64() < pc_p.ChosenProportion {
					// This assumes our value is copied.
					newBody[i][k] = n1[i][k]
				} else {
					newBody[i][k] = n2[i][k]
				}
			}
		}
		nn[j] = neural.ModularNetwork{
			Body:      newBody,
			Activator: nn[index1].Activator,
		}
	}

	return nn
}
