package eda

import (
	"fmt"
	"testing"

	"bitbucket.org/StephenPatrick/goevo/eda/fitness"
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
