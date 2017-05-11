package eda

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"bitbucket.org/StephenPatrick/goevo/mut"

	"github.com/200sc/go-dist/floatrange"
	"github.com/stretchr/testify/assert"
)

func Seed() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func TestOneMaxPBIL(t *testing.T) {
	fmt.Println("OneMaxPBIL")
	Seed()
	length := 1000.0
	model, err := Loop(PBILModel,
		BenchTest,
		FitnessFunc(OnemaxABS),
		Length(int(length)),
		LearningRate(0.2),
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

func TestFourPeaksPBIL(t *testing.T) {
	fmt.Println("FourPeakPBIL")
	Seed()
	length := 100.0
	model, err := Loop(PBILModel,
		BenchTest,
		FitnessFunc(FourPeaks(int(length/10))),
		Length(int(length)),
		LearningRate(0.2),
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
