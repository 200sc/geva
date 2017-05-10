package eda

import (
	"fmt"

	"bitbucket.org/StephenPatrick/goevo/env"
	"bitbucket.org/StephenPatrick/goevo/evoerr"
)

type UMDA struct {
	Base
}

func (umda *UMDA) Continue() bool {
	fitness := umda.fitness(umda)
	fmt.Println(fitness, umda.goalFitness)
	return fitness > umda.goalFitness
}

func (umda *UMDA) Adjust() Model {
	// Classically the UMDA begins iterations as a population
	// We're rotating this algorithm so that it begins as distribution
	// and generates two populations within the adjust function

	// Generate a population of size samples by sampling umda

	// Generate fitnesses for the population

	// Select a sub population of size learningSamples from samples

	// Sum over all values in the learningSamples for each index and
	// divide the resulting sums by the total number of samples
	// to obtain a new umda

	// return that

	return umda
}

func UMDAModel(opts ...Option) (Model, error) {
	umda := new(UMDA)
	umda.Base = DefaultBase()
	for _, opt := range opts {
		opt(umda)
	}
	if umda.length <= 0 {
		return nil, evoerr.InvalidLengthError{}
	}
	umda.F = env.NewF(umda.length, umda.baseValue)
	if umda.randomize {
		umda.F.RandomizeSingle(0.0, 1.0)
	}
	return umda, nil
}
