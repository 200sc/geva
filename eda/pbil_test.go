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

func TestOneMax(t *testing.T) {
	Seed()
	rng := floatrange.NewLinear(0.0, 1.0)
	model, err := Loop(PBILModel, 100,
		FitnessFunc(OnemaxABS),
		GoalFitness(4),
		Length(1000),
		BaseValue(0.5),
		Randomize(true),
		LearningRate(0.20),
		MutationRate(0.03),
		FMutator( //mut.Or(
			mut.And(
				mut.LinearRange(0.10),
				rng.EnforceRange),
			//mut.DropOut(0.5),
			//0.999,
		),
	)
	assert.Nil(t, err)
	fmt.Println(model.ToEnv())
}
