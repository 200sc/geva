package dev

import (
	"github.com/200sc/geva/cross"
	"github.com/200sc/geva/mut/frange"
	"github.com/200sc/go-dist/floatrange"
	"github.com/200sc/go-dist/intrange"
	"github.com/oakmound/oak/alg"
)

type Crossover interface {
	Crossover(a, b *Base) *Base
}

type ActionTypeCrossover interface {
	Crossover(a, b *Base) ([]ActionType, []float64, []floatrange.Range)
}

type ActionModCrossover struct {
	TypeLengthMod intrange.Range
	TypeWeightMod floatrange.Range
	TypeWeightDef floatrange.Range
	StrengthMod   frange.Mutator
}

func (amc *ActionModCrossover) Crossover(a, b *Base) ([]ActionType, []float64, []floatrange.Range) {
	actionTypeCount := len(a.ActionTypes) + len(b.ActionTypes)/2
	actionTypeCount += amc.TypeLengthMod.Poll()
	aTypeCount := actionTypeCount / 2
	bTypeCount := actionTypeCount - aTypeCount // In case of odd actionTypeCount

	actionTypes := make([]ActionType, actionTypeCount)
	actionStrengths := make([]floatrange.Range, actionTypeCount)
	actionTypeWeights := make([]float64, actionTypeCount)

	weights := make([]float64, len(a.ActionTypes))
	for i := 0; i < len(weights); i++ {
		weights[i] = 1.0
	}
	aChosen := alg.ChooseX(weights, aTypeCount)

	weights = make([]float64, len(b.ActionTypes))
	for i := 0; i < len(weights); i++ {
		weights[i] = 1.0
	}
	bChosen := alg.ChooseX(weights, bTypeCount)

	j := 0
	for i := 0; i < len(aChosen); i++ {
		k := aChosen[i]
		actionTypes[j] = a.ActionTypes[k]
		actionStrengths[j] = amc.StrengthMod(a.ActionStrengths[k])
		actionTypeWeights[j] = a.ActionTypeWeights[k] + amc.TypeWeightMod.Poll()
		j++
	}
	for i := 0; i < len(bChosen); i++ {
		k := bChosen[i]
		actionTypes[j] = b.ActionTypes[k]
		actionStrengths[j] = amc.StrengthMod(b.ActionStrengths[k])
		actionTypeWeights[j] = b.ActionTypeWeights[k] + amc.TypeWeightMod.Poll()
		j++
	}
	for i := 0; i < len(actionTypes); i++ {
		if actionTypeWeights[i] < 0 {
			actionTypeWeights[i] = amc.TypeWeightDef.Poll()
		}
	}

	return actionTypes, actionTypeWeights, actionStrengths
}

type LinearDevCrossover struct {
	ActionTypeCrossover
}

func (ldc *LinearDevCrossover) Crossover(a, b *Base) *Base {
	c := new(Base)
	c.InitSize = cross.LinearIntRange(a.InitSize, b.InitSize)
	c.GoalSize = cross.LinearIntRange(a.GoalSize, b.GoalSize)
	c.GoalDistance = cross.LinearIntRange(a.GoalDistance, b.GoalDistance)
	c.EnvSize = cross.LinearIntRange(a.EnvSize, b.EnvSize)
	c.EnvVal = cross.LinearFloatRange(a.EnvVal, b.EnvVal)
	c.ActionCount = cross.LinearIntRange(a.ActionCount, b.ActionCount)

	c.ActionTypes, c.ActionTypeWeights, c.ActionStrengths = ldc.ActionTypeCrossover.Crossover(a, b)

	return c
}
