package dev

import (
	"math/rand"

	"github.com/200sc/go-dist/floatrange"
	"github.com/200sc/go-dist/intrange"
	"github.com/oakmound/oak/alg"
)

type Creator interface {
	NewDev() Dev
}

type LinearCreator struct {
	InitSizeBottom       floatrange.Range
	InitSizeTop          floatrange.Range
	GoalSizeBottom       floatrange.Range
	GoalSizeTop          floatrange.Range
	GoalDistanceBottom   floatrange.Range
	GoalDistanceTop      floatrange.Range
	EnvSizeBottom        floatrange.Range
	EnvSizeTop           floatrange.Range
	EnvValBottom         floatrange.Range
	EnvValTop            floatrange.Range
	ActionCountBottom    floatrange.Range
	ActionCountTop       floatrange.Range
	ActionTypeCount      intrange.Range
	ActionTypeChoices    []ActionType
	ActionStrengthBottom floatrange.Range
	ActionStrengthTop    floatrange.Range
	CrossoverOptions     []Crossover
	MutationOptions      []DevMutation
}

func (ldc *LinearCreator) NewDev() Dev {
	return newDevFromRanges(ldc, floatrange.NewLinear)
}

// SpreadDevCreators just use the same fields as Linear
// but with bottom=base and top=spread
type SpreadCreator struct {
	*LinearCreator
}

func (bsdc *SpreadCreator) NewDev() Dev {
	return newDevFromRanges(bsdc.LinearCreator, floatrange.NewSpread)
}

func newDevFromRanges(ldc *LinearCreator, rngFn func(float64, float64) floatrange.Range) Dev {
	actTypeCt := ldc.ActionTypeCount.Poll()

	actTypes := make([]ActionType, actTypeCt)
	actStrengths := make([]floatrange.Range, actTypeCt)

	weights := make([]float64, len(ldc.ActionTypeChoices))
	for i := 0; i < len(weights); i++ {
		weights[i] = 1.0
	}
	chosen := alg.ChooseX(weights, actTypeCt)
	for i := 0; i < actTypeCt; i++ {
		actTypes[i] = ldc.ActionTypeChoices[chosen[i]]
		actStrengths[i] = rngFn(ldc.ActionStrengthBottom.Poll(), ldc.ActionStrengthTop.Poll())
	}

	return &Base{
		InitSize:        roundRange(rngFn(ldc.InitSizeBottom.Poll(), ldc.InitSizeTop.Poll())),
		GoalSize:        roundRange(rngFn(ldc.GoalSizeBottom.Poll(), ldc.GoalSizeTop.Poll())),
		GoalDistance:    roundRange(rngFn(ldc.GoalDistanceBottom.Poll(), ldc.GoalDistanceTop.Poll())),
		EnvSize:         roundRange(rngFn(ldc.EnvSizeBottom.Poll(), ldc.EnvSizeTop.Poll())),
		EnvVal:          rngFn(ldc.EnvValBottom.Poll(), ldc.EnvValTop.Poll()),
		ActionCount:     roundRange(rngFn(ldc.ActionCountBottom.Poll(), ldc.ActionCountTop.Poll())),
		ActionTypes:     actTypes,
		ActionStrengths: actStrengths,
		Cross:           ldc.CrossoverOptions[rand.Intn(len(ldc.CrossoverOptions))],
		DevMutation:     ldc.MutationOptions[rand.Intn(len(ldc.MutationOptions))],
	}
}

func roundRange(f floatrange.Range) intrange.Range {
	min := f.Percentile(0)
	max := f.Percentile(1)
	return intrange.NewLinear(alg.RoundF64(min), alg.RoundF64(max))
}
