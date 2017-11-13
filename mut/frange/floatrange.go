package frange

import (
	"math/rand"

	"github.com/200sc/go-dist/floatrange"
)

type Mutator func(floatrange.Range) floatrange.Range

// DropOut resets a range to some default
func DropOut(setTo floatrange.Range) Mutator {
	return func(f floatrange.Range) floatrange.Range {
		return setTo
	}
}

// Scale scales the input range by s
func Scale(s float64) Mutator {
	return func(f floatrange.Range) floatrange.Range {
		return f.Mult(s)
	}
}

// Div divides the input range by d
func Div(d float64) Mutator {
	return func(f floatrange.Range) floatrange.Range {
		return f.Mult(1 / d)
	}
}

// Add combines a and f through addition
func Add(a float64) Mutator {
	return func(f floatrange.Range) floatrange.Range {
		min := f.Percentile(0)
		max := f.Percentile(1)
		return floatrange.NewLinear(min+a, max+a)
	}
}

func EnforceMin(mn float64) Mutator {
	return func(f floatrange.Range) floatrange.Range {
		min := f.Percentile(0)
		max := f.Percentile(1)
		if min < mn {
			if max < mn {
				max = mn + 1
			}
			return floatrange.NewLinear(mn, max)
		}
		if max < mn {
			return floatrange.NewLinear(min, min+1)
		}
		return f
	}
}

// None performs no mutation on f
func None() Mutator {
	return func(f floatrange.Range) floatrange.Range {
		return f
	}
}

// And performs two range mutations in order
func And(a, b Mutator) Mutator {
	return func(f floatrange.Range) floatrange.Range {
		return b(a(f))
	}
}

// Or will perform a at chance aChance, and otherwise will
// perform b.
func Or(a, b Mutator, aChance float64) Mutator {
	return func(f floatrange.Range) floatrange.Range {
		if rand.Float64() < aChance {
			return a(f)
		}
		return b(f)
	}
}
