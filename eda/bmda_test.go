package eda

import (
	"fmt"
	"testing"

	"bitbucket.org/StephenPatrick/goevo/eda/fitness"
	"bitbucket.org/StephenPatrick/goevo/selection"
)

func TestOneMaxBMDA(t *testing.T) {
	fmt.Println("OneMaxBMDA")
	length := 100.0
	Loop(BMDAModel,
		BenchTest,
		FitnessFunc(fitness.OnemaxABS),
		Length(int(length)),
		SelectionMethod(selection.DeterministicTournament{5, 1}),
	)
}

func TestFourPeaksBMDA(t *testing.T) {
	fmt.Println("FourPeakBMDA")
	length := 100.0
	Loop(BMDAModel,
		BenchTest,
		FitnessFunc(fitness.FourPeaks(int(length/10))),
		Length(int(length)),
		SelectionMethod(selection.DeterministicTournament{5, 1}),
		MutationRate(.35),
	)
}

func TestQuadraticBMDA(t *testing.T) {
	fmt.Println("QuadraticBMDA")
	length := 100.0
	Loop(BMDAModel,
		BenchTest,
		FitnessFunc(fitness.Quadratic),
		Length(int(length)),
		SelectionMethod(selection.DeterministicTournament{5, 1}),
		MutationRate(.25),
	)
}

// BMDA tests:
// Deceptive n^3
// NK
// we are avoiding both of these because we're just looking at bitstrings
