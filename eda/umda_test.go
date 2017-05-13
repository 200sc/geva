package eda

import (
	"fmt"
	"testing"

	"bitbucket.org/StephenPatrick/goevo/eda/fitness"
	"bitbucket.org/StephenPatrick/goevo/selection"
)

func TestOneMaxUMDA(t *testing.T) {
	fmt.Println("OneMaxUMDA")
	length := 1000.0
	Loop(UMDAModel,
		BenchTest,
		FitnessFunc(fitness.OnemaxABS),
		Length(int(length)),
		SelectionMethod(selection.DeterministicTournament{5, 1}),
		MutationRate(0.03),
	)
}
