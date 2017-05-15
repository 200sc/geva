package eda

import (
	"fmt"
	"testing"

	"bitbucket.org/StephenPatrick/goevo/eda/fitness"
)

func TestFourPeaksECGA(t *testing.T) {
	fmt.Println("FourPeakECGA")
	length := 100.0
	Loop(ECGAModel,
		BenchTest,
		FitnessFunc(fitness.FourPeaks(int(length/10))),
		Length(int(length)),
		LearningRate(0.3),
		MutationRate(0.01),
	)
}

func TestTrap3ECGA(t *testing.T) {
	fmt.Println("Trap3ECGA")
	length := 100.0
	Loop(ECGAModel,
		BenchTest,
		FitnessFunc(fitness.TrapABS(3)),
		Length(int(length)),
		LearningRate(0.2),
		MutationRate(0.001),
	)
}
