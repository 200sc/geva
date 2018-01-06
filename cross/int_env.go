package cross

import (
	"math"

	"github.com/200sc/geva/env"
)

type I interface {
	Crossover(*env.I, *env.I) *env.I
}

type IPoint struct {
	NumPoints int
}

func (ipc IPoint) Crossover(a *env.I, b *env.I) *env.I {
	points := RandomPoints(ipc.NumPoints)

	short := b
	if len(*a) < len(*b) {
		short = a
	}
	active := a
	inactive := b
	start := 0
	end := 0

	var c []*int

	// Populate our new empty individual
	// combining the two parent individuals
	// as according to the above split points
	for _, v := range points {
		end = int(math.Ceil(float64(len(*short)) * v))

		c = append(c, (*active)[start:end]...)

		active, inactive = inactive, active
		start = end
	}

	// Add the remaining elements from the last individual
	c2 := env.I(append(c, (*active)[start:]...))
	return &c2
}
