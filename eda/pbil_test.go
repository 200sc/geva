package eda

import (
	"fmt"
	"testing"

	"bitbucket.org/StephenPatrick/goevo/eda/fitness"

	"github.com/stretchr/testify/assert"
)

func TestOneMaxPBIL(t *testing.T) {
	fmt.Println("OneMaxPBIL")
	Seed()
	length := 1000.0
	model, err := Loop(PBILModel,
		BenchTest,
		FitnessFunc(fitness.OnemaxABS),
		Length(int(length)),
		LearningRate(0.5),
		MutationRate(0.03),
	)
	assert.Nil(t, err)
	assert.NotNil(t, model)
}

func TestFourPeaksPBIL(t *testing.T) {
	fmt.Println("FourPeakPBIL")
	Seed()
	length := 100.0
	model, err := Loop(PBILModel,
		BenchTest,
		FitnessFunc(fitness.FourPeaks(int(length/10))),
		Length(int(length)),
		LearningRate(0.2),
		MutationRate(0.03),
	)
	assert.Nil(t, err)
	assert.NotNil(t, model)
}
