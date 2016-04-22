package crossover

import (
	"goevo/neural"
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

func (ac AverageCrossover) Crossover(nn []neural.Network, populated int, pairs [][]int) []neural.Network {

	pairIndex := 0

	for j := populated; j < len(nn); j++ {

		n1 := nn[pairs[pairIndex][0]].Body
		n2 := nn[pairs[pairIndex][1]].Body

		newBody := n1.CopyStructure()
		// This assumes that each network has the same dimensions!
		// Some vector math libraries would be good here (and elsewhere of course)
		for i := 0; i < len(n1); i++ {
			for k := 0; k < len(n1[i]); k++ {
				newNeuron := make(neural.Neuron, len(n1[i][k]))
				for m := 0; m < len(n1[i][k]); m++ {
					newNeuron[m] = ((n1[i][k][m] * ac.WeightMod) + n2[i][k][m]) / (ac.WeightMod + 1)
				}
				newBody[i][k] = newNeuron
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
