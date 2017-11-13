package eda

import (
	"fmt"

	"github.com/200sc/geva/env"
	"github.com/200sc/geva/pop"
)

// EnvInd is a wrapper around an environment and serves as the basis
// for the UMDA population
type EnvInd struct {
	*env.F
}

// NewEnvInd initializes a EnvInd to be a sampling of an env
func NewEnvInd(e *env.F) *EnvInd {
	return &EnvInd{GetSample(e)}
}

// Fitness is a dummy function on a EnvInd
func (ei *EnvInd) Fitness(input, expected [][]float64) int {
	return 0
}

// Mutate is a NOP on a EnvInd
func (ei *EnvInd) Mutate() {}

// Crossover is  NOP on a EnvInd
func (ei *EnvInd) Crossover(other pop.Individual) pop.Individual {
	return ei
}

// CanCrossover always returns false for a EnvInd
func (ei *EnvInd) CanCrossover(other pop.Individual) bool {
	return false
}

// Print prints a EnvInd
func (ei *EnvInd) Print() {
	fmt.Println(ei.F)
}
