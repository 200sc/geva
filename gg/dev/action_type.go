package dev

import (
	"github.com/200sc/geva/mut"
	"github.com/200sc/go-dist/floatrange"
)

type ActionType func(strength float64) mut.FloatMutator

var (
	BaseActionTypes = []ActionType{
		mut.Add,
		mut.DropOut,
		mut.Div,
		mut.LinearRange,
		mut.Scale,
	}
	BaseActionWeights = []float64{
		0.5,
		0.3,
		0.05,
		0.1,
		0.1,
	}
	BaseActionStrength = []floatrange.Range{
		floatrange.NewLinear(-15, 15),
		floatrange.NewLinear(-30, 30),
		floatrange.NewLinear(5, 9),
		floatrange.NewLinear(1, 10),
		floatrange.NewLinear(-5, 5),
	}
)

type ActionMut func(mut.FloatMutator, *float64) func()

func actionMut(mt mut.FloatMutator, f *float64) func() {
	return func() {
		*f = mt(*f)
	}
}
