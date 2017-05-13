package eda

import (
	"fmt"
	"testing"

	"bitbucket.org/StephenPatrick/goevo/mut"
	"github.com/200sc/go-dist/floatrange"
	"github.com/stretchr/testify/assert"
)

func TestFourPeaksMIMIC(t *testing.T) {
	fmt.Println("FourPeakMIMIC")
	Seed()
	length := 100.0
	model, err := Loop(MIMICModel,
		BenchTest,
		FitnessFunc(FourPeaks(int(length/10))),
		Length(int(length)),
		LearningRate(0.07),
		MutationRate(.15),
		FMutator(
			mut.And(
				mut.Or(
					mut.Or(mut.Add(.1), mut.Add(-.1), .5),
					mut.DropOut(0.5), .99),
				EnforceRange(floatrange.NewLinear(0.0, 1.0))),
		),
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
		FitnessFunc(SixPeaks(int(length/10))),
		Length(int(length)),
		LearningRate(0.07),
		MutationRate(.15),
		FMutator(
			mut.And(
				mut.Or(
					mut.Or(mut.Add(.1), mut.Add(-.1), .5),
					mut.DropOut(0.5), .99),
				EnforceRange(floatrange.NewLinear(0.0, 1.0))),
		),
	)
	assert.Nil(t, err)
	assert.NotNil(t, model)
}
