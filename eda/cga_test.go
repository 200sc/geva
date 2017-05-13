package eda

import (
	"fmt"
	"testing"

	"bitbucket.org/StephenPatrick/goevo/eda/fitness"
	"bitbucket.org/StephenPatrick/goevo/mut"

	"github.com/200sc/go-dist/floatrange"
	"github.com/stretchr/testify/assert"
)

func TestOneMaxCGA(t *testing.T) {
	fmt.Println("OneMaxCGA")
	Seed()
	length := 1000.0
	model, err := Loop(CGAModel,
		BenchTest,
		FitnessFunc(fitness.OnemaxABS),
		Length(int(length)),
		LearningRate(0.1),
		MutationRate(.03),
	)
	assert.Nil(t, err)
	assert.NotNil(t, model)
}

func TestFourPeaksCGA(t *testing.T) {
	fmt.Println("FourPeakCGA")
	Seed()
	length := 100.0
	model, err := Loop(CGAModel,
		BenchTest,
		FitnessFunc(fitness.FourPeaks(int(length/10))),
		Length(int(length)),
		LearningRate(0.1),
		MutationRate(.03),
	)
	assert.Nil(t, err)
	assert.NotNil(t, model)
}

func EnforceRange(fr floatrange.Range) mut.FloatMutator {
	return fr.EnforceRange
}
