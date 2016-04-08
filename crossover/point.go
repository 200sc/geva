package crossover

import (
	"goevo/neural"
)

// Randomly determine numPoints points to stitch two networks
// together at. For each numPoints, a point in a similar position
// along both networks will be chosen to split at. This will be
// more consistent if neural networks cannot expand or reduce
// in size.
type PointCrossover struct {
	numPoints int
}

func (pc_p *PointCrossover) Crossover(nn []neural.Network) []neural.Network {
	return nn
}
