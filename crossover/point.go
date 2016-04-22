package crossover

import (
	//"fmt"
	"goevo/neural"
	"math"
	"math/rand"
)

// Randomly determine NumPoints points to stitch two networks
// together at. For each NumPoints, a point in a similar position
// along both networks will be chosen to split at. This will be
// more consistent if neural networks cannot expand or reduce
// in size.
type PointCrossover struct {
	NumPoints int
}

func (pc PointCrossover) Crossover(nn []neural.Network, populated int, pairs [][]int) []neural.Network {

	pairIndex := 0

	for j := populated; j < len(nn); j++ {

		n1 := nn[pairs[pairIndex][0]].Body
		n2 := nn[pairs[pairIndex][1]].Body

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

		newBody := make(neural.Body, 0)

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

		nn[j] = neural.Network{
			Body:      newBody,
			Activator: nn[pairs[pairIndex][0]].Activator,
		}

		pairIndex++
	}

	return nn
}
