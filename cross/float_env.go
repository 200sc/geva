package cross

import (
	"math"
	"math/rand"

	"github.com/200sc/geva/env"
)

type F interface {
	Crossover(a, b *env.F) *env.F
}

type FPointCrossover struct {
	NumPoints int
}

func (fpc FPointCrossover) Crossover(a, b *env.F) *env.F {
	inc := 1.0 / float64(fpc.NumPoints)
	points := make([]float64, fpc.NumPoints)

	// Generate a series of random points to split on
	i := 0
	for pointRange := 0.0; pointRange < 1.0; pointRange += inc {
		points[i] = (rand.Float64() / float64(fpc.NumPoints)) + pointRange
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

	var c []*float64

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
