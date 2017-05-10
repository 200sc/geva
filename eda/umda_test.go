package eda

import (
	"fmt"
	"testing"

	"bitbucket.org/StephenPatrick/goevo/selection"

	"github.com/stretchr/testify/assert"
)

func TestOneMaxUMDA(t *testing.T) {
	Seed()

	length := 1000.0
	model, err := Loop(UMDAModel,
		Samples(40),
		LearningSamples(10),
		//FitnessFunc(OnemaxChance),
		FitnessFunc(OnemaxABS),
		GoalFitness(4),
		Length(int(length)),
		BaseValue(0.5),
		Randomize(true),
		SelectionMethod(selection.DeterministicTournament{2, 1}),
	)
	assert.Nil(t, err)
	fmt.Println(model.ToEnv())
}
