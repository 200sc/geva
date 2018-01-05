package eda

import (
	"fmt"
	"testing"

	"github.com/200sc/geva/eda/fitness"
)

func TestOneMaxBMDA(t *testing.T) {
	fmt.Println("OneMaxBMDA")
	length := 100.0
	Loop(BMDAModel,
		BenchTest,
		FitnessFunc(fitness.OnemaxABS),
		Length(int(length)),
	)
}

func TestFourPeaksBMDA(t *testing.T) {
	fmt.Println("FourPeakBMDA")
	length := 100.0
	Loop(BMDAModel,
		BenchTest,
		FitnessFunc(fitness.FourPeaks(int(length/10))),
		Length(int(length)),
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
		MutationRate(.25),
	)
}

// BMDA tests:
// Deceptive n^3
// NK
// we are avoiding both of these because we're just looking at bitstrings
