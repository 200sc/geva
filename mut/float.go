package mut

import (
	"math/rand"

	"github.com/200sc/go-dist/floatrange"
)

// A FloatMutator is expected to manipulate float64s. Exactly how they do that
// is up to the individual function.
type FloatMutator func(float64) float64

// LinearRange mutates a float to be somewhere between f - range and f + range,
// linearly
func LinearRange(rnge float64) FloatMutator {
	return func(f float64) float64 {
		return floatrange.NewSpread(f, rnge).Poll()
	}
}

// DropOut resets a float to some default
func DropOut(setTo float64) FloatMutator {
	return func(f float64) float64 {
		return setTo
	}
}

// Scale scales the input float by s
func Scale(s float64) FloatMutator {
	return func(f float64) float64 {
		return f * s
	}
}

// None performs no mutation on f
func None() FloatMutator {
	return func(f float64) float64 {
		return f
	}
}

// And performs two float mutations in order
func And(a, b FloatMutator) FloatMutator {
	return func(f float64) float64 {
		return b(a(f))
	}
}

// Or will perform a at chance aChance, and otherwise will
// perform b.
func Or(a, b FloatMutator, aChance float64) FloatMutator {
	return func(f float64) float64 {
		if rand.Float64() < aChance {
			return a(f)
		}
		return b(f)
	}
}
