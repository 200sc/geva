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
	//lrng := floatrange.NewLinear(0.05, 1.0)
	length := 1000.0
	model, err := Loop(PBILModel,
		Samples(40),
		LearningSamples(3),
		//FitnessFunc(OnemaxChance),
		FitnessFunc(OnemaxABS),
		GoalFitness(4),
		Length(int(length)),
		BaseValue(0.5),
		Randomize(true),
		LearningRate(0.5),
		MutationRate(3.0/(length/10.0)),
		FMutator( //mut.Or(
			mut.And(
				mut.LinearRange(0.10),
				rng.EnforceRange),
			//mut.DropOut(0.5),
			//0.999,
		),
		// LMutator(
		// 	mut.And(
		// 		mut.Scale(0.999),
		// 		lrng.EnforceRange),
		// ),
	)
	assert.Nil(t, err)
	fmt.Println(model.ToEnv())
}
