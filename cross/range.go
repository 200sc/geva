package cross

import (
	"github.com/200sc/go-dist/floatrange"
	"github.com/200sc/go-dist/intrange"
	"github.com/oakmound/oak/alg"
)

// FloatRange represents any function which can crossover
// floatrange.Range types
type FloatRange func(a, b floatrange.Range) floatrange.Range

// LinearFloatRange combines two floatranges into a linear range
// from the average of the minimum and maximum values of the inputs
func LinearFloatRange(a, b floatrange.Range) floatrange.Range {
	aMin := a.Percentile(0)
	aMax := a.Percentile(1)
	bMin := b.Percentile(0)
	bMax := b.Percentile(1)
	cMin := (aMin + bMin) / 2
	cMax := (aMax + bMax) / 2
	return floatrange.NewLinear(cMin, cMax)
}

// IntRange represents any function which can crossover
// intrange.Range types
type IntRange func(a, b intrange.Range) intrange.Range

// LinearIntRange combines two intranges into a linear range
// from the average of the minimum and maximum values of the inputs.
func LinearIntRange(a, b intrange.Range) intrange.Range {
	aMin := a.Percentile(0)
	aMax := a.Percentile(1)
	bMin := b.Percentile(0)
	bMax := b.Percentile(1)
	cMin := alg.RoundF64(float64(aMin+bMin) / 2)
	cMax := alg.RoundF64(float64(aMax+bMax) / 2)
	return intrange.NewLinear(cMin, cMax)
}
