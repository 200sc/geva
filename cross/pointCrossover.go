package cross

import (
	"math"
	"math/rand"
)

// RandomPoints returns a slice of numPoints float64s,
// randomly distributed from 0.0 to 1.0, but limited
// so that each point falls somewhere in the 1st, 2nd,
// etc. portion of the slice as it would be split in
// to as many parts as there are points.
//
// e.g. 2 points
// out[0] will be between 0.0 and 0.5
// out[1] will be between 0.5 and 1.0wdwdwwwwwwwww
func RandomPoints(numPoints int) []float64 {
	inc := 1.0 / float64(numPoints)
	points := make([]float64, numPoints)
	i := 0
	for pointRange := 0.0; pointRange < 1.0; pointRange += inc {
		points[i] = (rand.Float64() / float64(numPoints)) + pointRange
		i++
	}
	return points
}

// PointCrossover will convert two slices of interfaces into a new slice
// composed of elements from each slice. This composition will contain
// long strings of components from each input slice, at the same indices
// they were at in the original. The number of uninterrupted strings of
// elements from one slice is noted by numPoints + 1, where numPoints 
// is the number of points the slices are split at. 
func PointCrossover(a []interface{}, b []interface{}, numPoints int) []interface{} {

	points := RandomPoints(numPoints)
	short := b
	if len(a) < len(b) {
		short = a
	}
	active := a
	inactive := b
	start := 0
	var c []interface{}

	// Populate our new empty individual
	// combining the two parent individuals
	// as according to the above split points
	for _, v := range points {
		end := int(math.Ceil(float64(len(short)) * v))

		c = append(c, active[start:end]...)

		active, inactive = inactive, active
		start = end
	}

	// Add the remaining elements from the last individual
	return append(c, active[start:]...)
}
