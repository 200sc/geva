package neural

import (
	"math"
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

func (ac AverageCrossover) Crossover(a, b *Network) *Network {
	n1 := a.Body
	n2 := b.Body

	newBody := n1.CopyStructure()
	// This assumes that each network has the same dimensions!
	// Some vector math libraries would be good here (and elsewhere of course)
	for i := 0; i < len(n1); i++ {
		for k := 0; k < len(n1[i]); k++ {
			newNeuron := make(Neuron, len(n1[i][k]))
			for m := 0; m < len(n1[i][k]); m++ {
				newNeuron[m] = ((n1[i][k][m] * ac.WeightMod) + n2[i][k][m]) / (ac.WeightMod + 1)
			}
			newBody[i][k] = newNeuron
		}
	}
	return &Network{
		Body:      newBody,
		Activator: a.Activator,
	}
}

// Randomly determine NumPoints points to stitch two networks
// together at. For each NumPoints, a point in a similar position
// along both networks will be chosen to split at. This will be
// more consistent if neural networks cannot expand or reduce
// in size.
type PointCrossover struct {
	NumPoints int
}

func (pc PointCrossover) Crossover(a, b *Network) *Network {
	n1 := a.Body
	n2 := b.Body

	// Inc here is the value we use to ensure
	// distance between points-- each random
	// point is given an equal portion to be fit into.
	// This means a higher number of points will be
	// more uniform for smaller networks. numPoints
	// higher than 5 is not recommended, but hey
	// what do I know
	inc := 1.0 / float64(pc.NumPoints)
	points := make([]float64, pc.NumPoints)

	// Generate a series of random points to split on
	i := 0
	for pointRange := 0.0; pointRange < 1.0; pointRange += inc {
		r := (rand.Float64() / float64(pc.NumPoints)) + pointRange
		points[i] = r
		i++
	}

	activeNetwork := n1
	activeIndex := 1
	start := 0
	end := 0

	newBody := make(Body, 0)

	// Populate our new empty network by
	// combining the two parent networks
	// as according to the above split points
	for _, v := range points {
		end = int(math.Ceil(float64(len(activeNetwork)) * v))
		newBody = append(newBody, activeNetwork[start:end]...)

		if activeIndex == 1 {
			activeNetwork = n2
			activeIndex = 2
		} else {
			activeNetwork = n1
			activeIndex = 1
		}
		start = end
	}
	// Add on the remaining columns from the last network.
	end = int(math.Ceil(float64(len(activeNetwork)) * points[len(points)-1]))
	newBody = append(newBody, activeNetwork[end:]...)

	return &Network{
		Body:      newBody,
		Activator: a.Activator,
	}
}

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

func (uc *UniformCrossover) Crossover(a, b *Network) *Network {
	n1 := a.Body
	n2 := b.Body
	newBody := n1.CopyStructure()
	// This assumes that each network has the same dimensions!
	for i := 0; i < len(n1); i++ {
		for k := 0; k < len(n1[i]); k++ {
			if rand.Float64() < uc.ChosenProportion {
				// This assumes our value is copied and is not just a pointer.
				newBody[i][k] = n1[i][k]
			} else {
				newBody[i][k] = n2[i][k]
			}
		}
	}
	return &Network{
		Body:      newBody,
		Activator: a.Activator,
	}
}
