package eda

import (
	"fmt"
	"testing"

	"bitbucket.org/StephenPatrick/goevo/mut"

	"github.com/200sc/go-dist/floatrange"
	"github.com/stretchr/testify/assert"
)

func TestOneMaxCGA(t *testing.T) {
	Seed()
	length := 1000.0
	model, err := Loop(CGAModel,
		//FitnessFunc(OnemaxChance),
		FitnessFunc(OnemaxABS),
		GoalFitness(4),
		Length(int(length)),
		BaseValue(0.5),
		//Randomize(true),
		LearningRate(0.1),
		MutationRate(3.0/(length/10.0)),
		FMutator( //mut.Or(
			mut.And(
				mut.Or(mut.Add(.1), mut.Add(-.1), .5),
				EnforceRange(floatrange.NewLinear(0.0, 1.0))),
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

func EnforceRange(fr floatrange.Range) mut.FloatMutator {
	return fr.EnforceRange
}
