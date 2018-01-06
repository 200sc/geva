package cross

import (
	"math"

	"github.com/200sc/geva/env"
)

type F interface {
	Crossover(a, b *env.F) *env.F
}

type FPoint struct {
	NumPoints int
}

func (fpc FPoint) Crossover(a, b *env.F) *env.F {
	points := RandomPoints(fpc.NumPoints)

	short := b
	if len(*a) < len(*b) {
		short = a
	}
	active := a
	inactive := b
	start := 0
	end := 0

	var c []*float64

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
	c2 := env.F(append(c, (*active)[start:]...))
	return &c2
}

type FAverageCrossover struct {
	AWeight float64
}

func (fac FAverageCrossover) Crossover(a, b *env.F) *env.F {
	a2 := a.Copy().Mult(fac.AWeight / .5)
	b2 := b.Copy().Mult(.5 / fac.AWeight)
	return a2.AddF(b2).Divide(2)
}
