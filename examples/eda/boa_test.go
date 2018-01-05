package eda

import (
	"fmt"
	"testing"

	"github.com/200sc/geva/eda/fitness"
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
