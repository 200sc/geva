package crossover

import (
	"goevo/neural"
	"math/rand"
)

// For every neuron in the two networks, take the weights
// that neuron has and average them for a new network.
// They'll be averaged by ((weight1 * weightMod) + weight2) / (weightMod + 1)
type AverageCrossover struct {
	// This weight is applied to all weights in the first
	// network selected, before the average of the networks
	// is calculated. A weight more distant from 1 will
	// swing the averaged networks toward more closely
	// emulating one network or the other. Cannot be negative.
	//
	// This might need to be modified into two weightMods
	// if crossover pairings are determined non-randomly
	WeightMod float64
}

func (ac AverageCrossover) Crossover(nn []neural.ModularNetwork, populated int) []neural.ModularNetwork {

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
		// Some vector math libraries would be good here (and elsewhere of course)
		for i := 0; i < len(n1); i++ {
			for k := 0; k < len(n1[i]); k++ {
				newNeuron := make(neural.ModularNeuron, len(n1[i][k]))
				for m := 0; m < len(n1[i][k]); m++ {
					newNeuron[m] = ((n1[i][k][m] * ac.WeightMod) + n2[i][k][m]) / (ac.WeightMod + 1)
				}
				newBody[i][k] = newNeuron
			}
		}
		nn[j] = neural.ModularNetwork{
			Body:      newBody,
			Activator: nn[index1].Activator,
		}
	}

	return nn
}
