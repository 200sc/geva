package cross

import (
	"math"
	"math/rand"

	"github.com/200sc/geva/env"
)

type I interface {
	Crossover(*env.I, *env.I) *env.I
}

type IPointCrossover struct {
	NumPoints int
}

func (ipc IPointCrossover) Crossover(a *env.I, b *env.I) *env.I {
	inc := 1.0 / float64(ipc.NumPoints)
	points := make([]float64, ipc.NumPoints)

	// Generate a series of random points to split on
	i := 0
	for pointRange := 0.0; pointRange < 1.0; pointRange += inc {
		points[i] = (rand.Float64() / float64(ipc.NumPoints)) + pointRange
		i++
	}

	short := b
	if len(*a) < len(*b) {
		short = a
	}
	activeFlag := false
	active := a
	start := 0
	end := 0

	var c []*int

	// Populate our new empty individual
	// combining the two parent individuals
	// as according to the above split points
	for _, v := range points {
		end = int(math.Ceil(float64(len(*short)) * v))

		c = append(c, (*active)[start:end]...)

		if !activeFlag {
			active = b
		} else {
			active = a
		}
		activeFlag = !activeFlag
		start = end
	}

	// Add the remaining elements from the last individual
	c2 := env.I(append(c, (*active)[start:]...))
	return &c2
}
