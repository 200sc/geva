package eda

import (
	"fmt"
	"testing"

	"bitbucket.org/StephenPatrick/goevo/eda/fitness"
	"bitbucket.org/StephenPatrick/goevo/mut"
	"bitbucket.org/StephenPatrick/goevo/selection"
	"github.com/200sc/go-dist/floatrange"
	"github.com/stretchr/testify/assert"
)

func TestOneMaxBMDA(t *testing.T) {
	fmt.Println("OneMaxBMDA")
	Seed()
	length := 100.0
	model, err := Loop(BMDAModel,
		BenchTest,
		FitnessFunc(fitness.OnemaxABS),
		Length(int(length)),
		SelectionMethod(selection.DeterministicTournament{2, 1}),
		MutationRate(.15),
		FMutator(
			mut.And(
				mut.Or(
					mut.Or(mut.Add(.1), mut.Add(-.1), .5),
					mut.DropOut(0.5), .99),
				EnforceRange(floatrange.NewLinear(0.0, 1.0))),
		),
	)
	assert.Nil(t, err)
	assert.NotNil(t, model)
}

func TestFourPeaksBMDA(t *testing.T) {
	fmt.Println("FourPeakBMDA")
	Seed()
	length := 100.0
	model, err := Loop(BMDAModel,
		BenchTest,
		FitnessFunc(fitness.FourPeaks(int(length/10))),
		Length(int(length)),
		SelectionMethod(selection.DeterministicTournament{4, 1}),
		MutationRate(.25),
		FMutator(
			mut.And(
				mut.Or(
					mut.Or(mut.Add(.1), mut.Add(-.1), .5),
					mut.DropOut(0.5), .99),
				EnforceRange(floatrange.NewLinear(0.0, 1.0))),
		),
	)
	assert.Nil(t, err)
	assert.NotNil(t, model)
}

func TestQuadraticBMDA(t *testing.T) {
	fmt.Println("QuadraticBMDA")
	Seed()
	length := 100.0
	model, err := Loop(BMDAModel,
		BenchTest,
		FitnessFunc(fitness.Quadratic),
		Length(int(length)),
		SelectionMethod(selection.DeterministicTournament{4, 1}),
		MutationRate(.25),
		FMutator(
			mut.And(
				mut.Or(
					mut.Or(mut.Add(.1), mut.Add(-.1), .5),
					mut.DropOut(0.5), .99),
				EnforceRange(floatrange.NewLinear(0.0, 1.0))),
		),
	)
	assert.Nil(t, err)
	assert.NotNil(t, model)
}

// BMDA tests:
// Deceptive n^3
// NK
// we are avoiding both of these because we're just looking at bitstrings
