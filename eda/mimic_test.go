package eda

import (
	"fmt"
	"testing"

	"bitbucket.org/StephenPatrick/goevo/eda/fitness"
	"github.com/stretchr/testify/assert"
)

func TestFourPeaksMIMIC(t *testing.T) {
	fmt.Println("FourPeakMIMIC")
	Seed()
	length := 100.0
	model, err := Loop(MIMICModel,
		BenchTest,
		FitnessFunc(fitness.FourPeaks(int(length/10))),
		Length(int(length)),
		LearningRate(0.07),
	)
	assert.Nil(t, err)
	assert.NotNil(t, model)
}

func TestSixPeaksMIMIC(t *testing.T) {
	fmt.Println("SixPeakMIMIC")
	Seed()
	length := 100.0
	model, err := Loop(MIMICModel,
		BenchTest,
		FitnessFunc(fitness.SixPeaks(int(length/10))),
		Length(int(length)),
		LearningRate(0.07),
	)
	assert.Nil(t, err)
	assert.NotNil(t, model)
}
