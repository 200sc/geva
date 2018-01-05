package eda

import (
	"fmt"
	"testing"

	"github.com/200sc/geva/eda/fitness"
)

func TestFourPeaksMIMIC(t *testing.T) {
	fmt.Println("FourPeakMIMIC")
	length := 100.0
	Loop(MIMICModel,
		BenchTest,
		FitnessFunc(fitness.FourPeaks(int(length/10))),
		Length(int(length)),
		LearningRate(0.07),
	)
}

func TestSixPeaksMIMIC(t *testing.T) {
	fmt.Println("SixPeakMIMIC")
	length := 100.0
	Loop(MIMICModel,
		BenchTest,
		FitnessFunc(fitness.SixPeaks(int(length/10))),
		Length(int(length)),
		LearningRate(0.07),
	)
}
