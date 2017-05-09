package mut

import "github.com/200sc/go-dist/floatrange"

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

// None performs no mutation on f
func None() FloatMutator {
	return func(f float64) float64 {
		return f
	}
}
