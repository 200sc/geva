package eda

import (
	"fmt"
	"testing"

	"github.com/200sc/geva/eda/fitness"
	"github.com/200sc/geva/mut"
	"github.com/200sc/go-dist/floatrange"
)

func TestOneMaxSHCLVND(t *testing.T) {
	fmt.Println("OneMaxSCHLVND")
	length := 100
	Loop(SHCLVNDModel,
		BenchTest,
		FitnessFunc(fitness.OnemaxABS),
		Length(int(length)),
		LearningRate(0.05),
		LMutator(
			mut.And(
				mut.Scale(0.997),
				EnforceRange(floatrange.NewLinear(0.001, 1.0))),
		),
	)
}
