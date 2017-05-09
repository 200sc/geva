package eda

import (
	"fmt"
	"testing"

	"bitbucket.org/StephenPatrick/goevo/mut"

	"github.com/stretchr/testify/assert"
)

func TestOneMax(t *testing.T) {
	model, err := Loop(PBILModel, 5,
		FitnessFunc(OnemaxABS),
		GoalFitness(3),
		Length(100),
		BaseValue(0.5),
		LearningRate(0.04),
		MutationRate(0.05),
		FMutator(mut.LinearRange(0.05)),
	)
	assert.Nil(t, err)
	fmt.Println(model.ToEnv())
}
