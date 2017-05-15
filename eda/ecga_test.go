package eda

import (
	"fmt"
	"testing"

	"bitbucket.org/StephenPatrick/goevo/eda/fitness"
)

func TestFourPeaksECGA(t *testing.T) {
	fmt.Println("FourPeakCGA")
	length := 100.0
	Loop(ECGAModel,
		BenchTest,
		FitnessFunc(fitness.FourPeaks(int(length/10))),
		Length(int(length)),
		LearningRate(0.1),
		MutationRate(.03),
	)
}
