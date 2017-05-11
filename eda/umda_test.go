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
	fmt.Println("OneMaxUMDA")
	Seed()
	length := 1000.0
	model, err := Loop(UMDAModel,
		BenchTest,
		FitnessFunc(OnemaxABS),
		Length(int(length)),
		SelectionMethod(selection.DeterministicTournament{2, 1}),
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
