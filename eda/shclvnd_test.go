package eda

import (
	"fmt"
	"testing"

	"bitbucket.org/StephenPatrick/goevo/mut"
	"github.com/200sc/go-dist/floatrange"
	"github.com/stretchr/testify/assert"
)

func TestOneMaxSHCLVND(t *testing.T) {
	Seed()
	length := 100
	model, err := Loop(SHCLVNDModel,
		Samples(200),
		LearningSamples(3),
		//FitnessFunc(OnemaxChance),
		FitnessFunc(OnemaxABS),
		GoalFitness(0),
		Length(int(length)),
		BaseValue(0.5),
		//Randomize(true),
		LearningRate(0.05),
		// MutationRate(3.0/(length/10.0)),
		// FMutator( //mut.Or(
		// 	mut.And(
		// 		mut.Or(mut.Add(.1), mut.Add(-.1), .5),
		// 		EnforceRange(floatrange.NewLinear(0.0, 1.0))),
		// 	//mut.DropOut(0.5),
		// 	//0.999,
		// ),
		LMutator(
			mut.And(
				mut.Scale(0.997),
				EnforceRange(floatrange.NewLinear(0.001, 1.0))),
		),
	)
	assert.Nil(t, err)
	fmt.Println(model.ToEnv())
}
