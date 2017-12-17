package irange

import (
	"math/rand"

	"github.com/200sc/go-dist/intrange"
)

type Mutator func(intrange.Range) intrange.Range

// DropOut resets a range to some default
func DropOut(setTo intrange.Range) Mutator {
	return func(f intrange.Range) intrange.Range {
		return setTo
	}
}

// Scale scales the input range by s
func Scale(s float64) Mutator {
	return func(f intrange.Range) intrange.Range {
		return f.Mult(s)
	}
}

// Div divides the input range by d
func Div(d float64) Mutator {
	return func(f intrange.Range) intrange.Range {
		return f.Mult(1 / d)
	}
}

// Add combines a and f through addition
func Add(a int) Mutator {
	return func(f intrange.Range) intrange.Range {
		min := f.Percentile(0)
		max := f.Percentile(1)
		return intrange.NewLinear(min+a, max+a)
	}
}

func EnforceMax(mx int) Mutator {
	return func(f intrange.Range) intrange.Range {
		min := f.Percentile(0)
		max := f.Percentile(1)
		if max > mx {
			max = mx
			if max < min {
				min = max - 1
			}
			return intrange.NewLinear(min, max)
		}
		if mx < min {
			return intrange.NewLinear(mx-1, mx)
		}
		return f
	}
}

func EnforceMin(mn int) Mutator {
	return func(f intrange.Range) intrange.Range {
		min := f.Percentile(0)
		max := f.Percentile(1)
		if min < mn {
			min = mn
			if max < mn {
				max = mn + 1
			}
			return intrange.NewLinear(mn, max)
		}
		if max < mn {
			return intrange.NewLinear(min, min+1)
		}
		return f
	}
}

// None performs no mutation on f
func None() Mutator {
	return func(f intrange.Range) intrange.Range {
		return f
	}
}

// And performs two range mutations in order
func And(a, b Mutator) Mutator {
	return func(f intrange.Range) intrange.Range {
		return b(a(f))
	}
}

// Or will perform a at chance aChance, and otherwise will
// perform b.
func Or(a, b Mutator, aChance float64) Mutator {
	return func(f intrange.Range) intrange.Range {
		if rand.Float64() < aChance {
			return a(f)
		}
		return b(f)
	}
}

// OrAny will perform one of the given mutations, each
// with an equal chance. It's anticipated that chance * len(muts) >= 1,
// if this isn't the case occassionally the mutator will do nothing.
func OrAny(chance float64, muts ...Mutator) Mutator {
	return func(f intrange.Range) intrange.Range {
		r := rand.Float64()
		for _, m := range muts {
			if r < chance {
				return m(f)
			}
			r -= chance
		}
		return f
	}
}
