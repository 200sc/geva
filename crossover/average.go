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
	weightMod float64
}

func (pc_p *AverageCrossover) Crossover(nn []neural.Network, populated int) []neural.Network {
	return nn
}
