package eda

import (
	"fmt"
	"testing"

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
		FitnessFunc(OnemaxABS),
		Length(int(length)),
		LearningRate(0.1),
		MutationRate(3.0/(length/10.0)),
		FMutator(
			mut.And(
				mut.Or(mut.Add(.1), mut.Add(-.1), .5),
				EnforceRange(floatrange.NewLinear(0.0, 1.0))),
		),
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
		FitnessFunc(FourPeaks(int(length/10))),
		Length(int(length)),
		LearningRate(0.1),
		MutationRate(.03),
		FMutator(
			mut.And(
				mut.Or(mut.Add(.1), mut.Add(-.1), .5),
				EnforceRange(floatrange.NewLinear(0.0, 1.0))),
		),
	)
	assert.Nil(t, err)
	assert.NotNil(t, model)
}

func EnforceRange(fr floatrange.Range) mut.FloatMutator {
	return fr.EnforceRange
}
