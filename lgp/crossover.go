package lgp

import (
	"math"
	"math/rand"
)

type LGPCrossover interface {
	Crossover(a, b *LGP) *LGP
}

type PointCrossover struct {
	NumPoints int
}

func (pc PointCrossover) Crossover(a, b *LGP) *LGP {
	inc := 1.0 / float64(pc.NumPoints)
	points := make([]float64, pc.NumPoints)

	// Generate a series of random points to split on
	i := 0
	for pointRange := 0.0; pointRange < 1.0; pointRange += inc {
		points[i] = (rand.Float64() / float64(pc.NumPoints)) + pointRange
		i++
	}

	var short *LGP
	if len(a.Instructions) < len(b.Instructions) {
		short = a
	} else {
		short = b
	}
	active := a
	start := 0
	end := 0

	c := a.Copy()
	c.Instructions = []Instruction{}

	// Populate our new empty individual
	// combining the two parent individuals
	// as according to the above split points
	for _, v := range points {
		end = int(math.Ceil(float64(len(short.Instructions)) * v))

		c.Instructions = append(c.Instructions, active.Instructions[start:end]...)

		if active == a {
			active = b
		} else {
			active = a
		}
		start = end
	}

	// Add the remaining elements from the last individual
	c.Instructions = append(c.Instructions, active.Instructions[start:]...)

	if len(*a.Mem) > len(*b.Mem) {
		c.Mem = a.Mem.Copy()
	} else {
		c.Mem = b.Mem.Copy()
	}

	return c
}

// This looks like it doesn't work
type UniformCrossover struct {
	ChosenProportion float64
}

func (uc UniformCrossover) Crossover(a, b *LGP) *LGP {

	var c *LGP

	// This UniformCrossover will artificially lengthen everybody
	// Recommend pairing with high ShrinkMutation
	if len(a.Instructions) > len(b.Instructions) {
		c = b.Copy()
	} else {
		c = a.Copy()
	}

	i := 0
	for i < len(c.Instructions) {
		if rand.Float64() >= uc.ChosenProportion {
			c.Instructions[i] = b.Instructions[i]
		} else {
			c.Instructions[i] = a.Instructions[i]
		}
		i++
	}

	if len(a.Instructions) > len(b.Instructions) {
		c.Instructions = append(c.Instructions, a.Instructions[i:]...)
	} else {
		c.Instructions = append(c.Instructions, b.Instructions[i:]...)
	}

	if len(*a.Mem) > len(*b.Mem) {
		c.Mem = a.Mem.Copy()
	} else {
		c.Mem = b.Mem.Copy()
	}

	return c
}
