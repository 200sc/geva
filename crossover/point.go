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

func (pc PointCrossover) Crossover(nn []neural.Network, populated int) []neural.Network {

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

		inc := 1.0 / float64(pc.NumPoints)
		points := make([]float64, pc.NumPoints)

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

		nn[j] = activeNetwork.Make()

		for _, v := range points {
			end = int(math.Ceil(float64(activeNetwork.Length()) * v))
			nn[j] = nn[j].Append(activeNetwork.Slice(start, end))

			if activeIndex == 1 {
				activeNetwork = n2
				activeIndex = 2
			} else {
				activeNetwork = n1
				activeIndex = 1
			}
			start = end
		}
		end = int(math.Ceil(float64(activeNetwork.Length()) * points[len(points)-1]))
		nn[j] = nn[j].Append(activeNetwork.SliceToEnd(end))
	}

	return nn
}
