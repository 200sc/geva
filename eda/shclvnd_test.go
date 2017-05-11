package eda

import (
	"fmt"
	"testing"

	"bitbucket.org/StephenPatrick/goevo/mut"
	"github.com/200sc/go-dist/floatrange"
	"github.com/stretchr/testify/assert"
)

func TestOneMaxSHCLVND(t *testing.T) {
	fmt.Println("OneMaxSCHLVND")
	Seed()
	length := 1000
	model, err := Loop(SHCLVNDModel,
		BenchTest,
		FitnessFunc(OnemaxABS),
		Length(int(length)),
		LearningRate(0.05),
		// MutationRate(3.0/(length/10.0)),
		// FMutator(
		// 	mut.And(
		// 		mut.Or(mut.Add(.1), mut.Add(-.1), .5),
		// 		EnforceRange(floatrange.NewLinear(0.0, 1.0))),
		// ),
		LMutator(
			mut.And(
				mut.Scale(0.997),
				EnforceRange(floatrange.NewLinear(0.001, 1.0))),
		),
	)
	assert.Nil(t, err)
	assert.NotNil(t, model)
}
