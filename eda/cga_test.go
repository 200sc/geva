package eda

import (
	"fmt"
	"testing"

	"bitbucket.org/StephenPatrick/goevo/eda/fitness"
	"bitbucket.org/StephenPatrick/goevo/mut"

	"github.com/200sc/go-dist/floatrange"
)

func TestOneMaxCGA(t *testing.T) {
	fmt.Println("OneMaxCGA")
	length := 1000.0
	Loop(CGAModel,
		BenchTest,
		FitnessFunc(fitness.OnemaxABS),
		Length(int(length)),
		LearningRate(0.1),
		MutationRate(.03),
	)
}

func TestFourPeaksCGA(t *testing.T) {
	fmt.Println("FourPeakCGA")
	length := 100.0
	Loop(CGAModel,
		BenchTest,
		FitnessFunc(fitness.FourPeaks(int(length/10))),
		Length(int(length)),
		LearningRate(0.1),
		MutationRate(.03),
	)
}

func EnforceRange(fr floatrange.Range) mut.FloatMutator {
	return fr.EnforceRange
}
