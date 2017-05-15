package eda

import (
	"fmt"
	"testing"

	"bitbucket.org/StephenPatrick/goevo/eda/fitness"
)

func TestFourPeaksBOA(t *testing.T) {
	fmt.Println("FourPeakBOA")
	length := 100.0
	Loop(BOAModel,
		BenchTest,
		FitnessFunc(fitness.FourPeaks(int(length/10))),
		Length(int(length)),
		LearningRate(0.1),
		MutationRate(0.03),
	)
}
