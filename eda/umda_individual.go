package eda

import (
	"fmt"

	"bitbucket.org/StephenPatrick/goevo/env"
	"bitbucket.org/StephenPatrick/goevo/pop"
)

// UMDAIndividual is a wrapper around an environment and serves as the basis
// for the UMDA population
type UMDAIndividual struct {
	*env.F
}

// NewUMDAIndividual initializes a UMDAIndividual to be a sampling of an env
func NewUMDAIndividual(e *env.F) *UMDAIndividual {
	return &UMDAIndividual{GetSample(e)}
}

// Fitness is a dummy function on a UMDAIndividual
func (umdaI *UMDAIndividual) Fitness(input, expected [][]float64) int {
	return 0
}

// Mutate is a NOP on a UMDAIndividual
func (umdaI *UMDAIndividual) Mutate() {}

// Crossover is  NOP on a UMDAIndividual
func (umdaI *UMDAIndividual) Crossover(other pop.Individual) pop.Individual {
	return umdaI
}

// CanCrossover always returns false for a UMDAIndividual
func (umdaI *UMDAIndividual) CanCrossover(other pop.Individual) bool {
	return false
}

// Print prints a UMDAIndividual
func (umdaI *UMDAIndividual) Print() {
	fmt.Println("UMDAInvidual")
}
