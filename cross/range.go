package cross

import (
	"github.com/200sc/go-dist/floatrange"
	"github.com/200sc/go-dist/intrange"
	"github.com/oakmound/oak/alg"
)

type FloatRangeCrossover func(a, b floatrange.Range) floatrange.Range

func LinearFloatRange(a, b floatrange.Range) floatrange.Range {
	aMin := a.Percentile(0)
	aMax := a.Percentile(1)
	bMin := b.Percentile(0)
	bMax := b.Percentile(1)
	cMin := (aMin + bMin) / 2
	cMax := (aMax + bMax) / 2
	return floatrange.NewLinear(cMin, cMax)
}

func LinearIntRange(a, b intrange.Range) intrange.Range {
	aMin := a.Percentile(0)
	aMax := a.Percentile(1)
	bMin := b.Percentile(0)
	bMax := b.Percentile(1)
	cMin := alg.RoundF64(float64(aMin+bMin) / 2)
	cMax := alg.RoundF64(float64(aMax+bMax) / 2)
	return intrange.NewLinear(cMin, cMax)
}
