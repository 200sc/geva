package eda

import (
	"fmt"

	"bitbucket.org/StephenPatrick/goevo/env"
	"bitbucket.org/StephenPatrick/goevo/evoerr"
	"bitbucket.org/StephenPatrick/goevo/pop"
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

	// This is a bastardization of the evolutionary population model
	// because the evolutionary population model assumes you will follow its
	// rules and not replace its elements and call its NextGeneration function
	// and we are breaking all of its rules,
	// which implies there are problems with the evolutionary population model
	// in that there should be some lower tiered struct that can't do
	// NextGeneration but can be selected on and crossovered on, etc.

	// Generate a population of size samples by sampling umda
	p := new(pop.Population)
	p.Members = make([]pop.Individual, umda.samples)
	for i := 0; i < umda.samples; i++ {
		p.Members[i] = NewUMDAIndividual(umda.F)
	}

	// Generate fitnesses for the population
	p.Fitnesses = make([]int, umda.samples)
	for i := 0; i < umda.samples; i++ {
		umda.F = p.Members[i].(*UMDAIndividual).F
		p.Fitnesses[i] = umda.fitness(umda)
	}

	p.Size = umda.learningSamples

	fmt.Println(p.Fitnesses)

	// Select a sub population of size learningSamples from samples
	// interface problem: we need to be able to ensure that parentProportion here
	// leaves us with at least (ideally exactly) learningSamples members.
	subPop := umda.selection.Select(p)

	// Sum over all values in the learningSamples for each index and
	newenv := env.NewF(umda.length, 0.0)
	for _, ind := range subPop {
		for j, f := range *(ind.(*UMDAIndividual).F) {
			*(*newenv)[j] = *(*newenv)[j] + *f
		}
	}
	// divide the resulting sums by the total number of samples
	// to obtain a new umda
	newenv.Divide(float64(len(subPop)))
	umda.F = newenv
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

type UMDAIndividual struct {
	*env.F
}

func NewUMDAIndividual(e *env.F) *UMDAIndividual {
	sample := GetSample(e)
	return &UMDAIndividual{sample}
}

// Placeholder interface satisfying functions

func (umdaI *UMDAIndividual) Fitness(input, expected [][]float64) int {
	return 0
}
func (umdaI *UMDAIndividual) Mutate() {}

func (umdaI *UMDAIndividual) Crossover(other pop.Individual) pop.Individual {
	return umdaI
}

func (umdaI *UMDAIndividual) CanCrossover(other pop.Individual) bool {
	return false
}

func (umdaI *UMDAIndividual) Print() {
	fmt.Println("UMDAInvidual")
}
