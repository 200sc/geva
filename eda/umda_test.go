package eda

import (
	"fmt"
	"testing"

	"bitbucket.org/StephenPatrick/goevo/mut"
	"bitbucket.org/StephenPatrick/goevo/selection"

	"github.com/200sc/go-dist/floatrange"
	"github.com/stretchr/testify/assert"
)

func TestOneMaxUMDA(t *testing.T) {
	Seed()

	length := 1000.0
	model, err := Loop(UMDAModel,
		Samples(100),
		LearningSamples(30),
		//FitnessFunc(OnemaxChance),
		FitnessFunc(OnemaxABS),
		GoalFitness(4),
		Length(int(length)),
		BaseValue(0.5),
		//Randomize(true),
		SelectionMethod(selection.DeterministicTournament{2, 1}),
		MutationRate(3.0/(length/10.0)),
		FMutator( //mut.Or(
			mut.And(
				mut.Or(mut.Add(.1), mut.Add(-.1), .5),
				EnforceRange(floatrange.NewLinear(0.0, 1.0))),
			//mut.DropOut(0.5),
			//0.999,
		),
	)
	assert.Nil(t, err)
	fmt.Println(model.ToEnv())
}
