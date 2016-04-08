package crossover

import (
	"goevo/neural"
)

// Choose a bunch of random neurons from each network
// and make a new network out of them.
// I don't think this is a very good idea for neural
// networks, but we'll see.
type UniformCrossover struct {
	// This proportion of neurons are chosen
	// from the first network selected.
	// The remaining proporiton 1 - chosenProportion
	// come from the other network.
	// Cannot be negative.
	chosenProportion float64
}

func (pc_p *UniformCrossover) Crossover(nn []neural.Network) []neural.Network {
	return nn
}
