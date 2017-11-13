package eda

import (
	"fmt"
	"testing"

	"github.com/200sc/geva/eda/fitness"
)

func TestOneMaxUMDA(t *testing.T) {
	fmt.Println("OneMaxUMDA")
	length := 1000.0
	Loop(UMDAModel,
		BenchTest,
		FitnessFunc(fitness.OnemaxABS),
		Length(int(length)),
		MutationRate(0.03),
	)
}
