package eda

import (
	"bitbucket.org/StephenPatrick/goevo/env"
	"bitbucket.org/StephenPatrick/goevo/pop"
)

// UMDA represents the Univariate Marginal Distribution Algorithm
type UMDA struct {
	Base
}

// Adjust on a UMDA creates a population from sampling the umda's distribution,
// then selects from that population top fitness members, then sets its
// distribution to be the average of the selected members
func (umda *UMDA) Adjust() Model {
	// Classically the UMDA begins iterations as a population
	// We're rotating this algorithm so that it begins as distribution
	// and generates two populations within the adjust function

	p := umda.Pop()
	p.Size = umda.learningSamples

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
	umda.F.Mutate(umda.mutationRate, umda.fmutator)

	return umda
}

// UMDAModel initializes a UMDA EDA
func UMDAModel(opts ...Option) (Model, error) {
	var err error
	umda := new(UMDA)
	umda.Base, err = DefaultBase(opts...)
	return umda, err
}

func (umda *UMDA) Pop() *pop.Population {
	// This is a bastardization of the evolutionary population model
	// because the evolutionary population model assumes you will follow its
	// rules and don't replace its elements and call its NextGeneration function
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

	p.Size = umda.samples
	return p
}
