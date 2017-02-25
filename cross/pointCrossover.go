package cross

import (
	"math"
	"math/rand"
)

func PointCrossover(a []interface{}, b []interface{}, numPoints int) []interface{} {
	inc := 1.0 / float64(numPoints)
	points := make([]float64, numPoints)

	// Generate a series of random points to split on
	i := 0
	for pointRange := 0.0; pointRange < 1.0; pointRange += inc {
		points[i] = (rand.Float64() / float64(numPoints)) + pointRange
		i++
	}

	var short []interface{}
	if len(a) < len(b) {
		short = a
	} else {
		short = b
	}
	activeFlag := 0
	active := a
	start := 0
	end := 0

	var c []interface{}

	// Populate our new empty individual
	// combining the two parent individuals
	// as according to the above split points
	for _, v := range points {
		end = int(math.Ceil(float64(len(short)) * v))

		c = append(c, active[start:end]...)

		if activeFlag == 0 {
			active = b
			activeFlag = 1
		} else {
			active = a
			activeFlag = 0
		}
		start = end
	}

	// Add the remaining elements from the last individual
	return append(c, active[start:]...)
}
