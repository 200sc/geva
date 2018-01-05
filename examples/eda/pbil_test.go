package eda

import (
	"fmt"
	"testing"

	"github.com/200sc/geva/eda/fitness"
)

func TestOneMaxPBIL(t *testing.T) {
	fmt.Println("OneMaxPBIL")
	length := 1000.0
	Loop(PBILModel,
		BenchTest,
		FitnessFunc(fitness.OnemaxABS),
		Length(int(length)),
		LearningRate(0.5),
		MutationRate(0.03),
	)
}

func TestAlternatingPBIL(t *testing.T) {
	fmt.Println("AlternatingPBIL")
	length := 1000.0
	Loop(PBILModel,
		BenchTest,
		FitnessFunc(fitness.AlternatingABS),
		Length(int(length)),
		LearningRate(0.5),
		MutationRate(0.03),
	)
}

func TestFourPeaksPBIL(t *testing.T) {
	fmt.Println("FourPeakPBIL")
	length := 100.0
	Loop(PBILModel,
		BenchTest,
		FitnessFunc(fitness.FourPeaks(int(length/10))),
		Length(int(length)),
		LearningRate(0.2),
		MutationRate(0.03),
	)
}
