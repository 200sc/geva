package eda

import (
	"bitbucket.org/StephenPatrick/goevo/env"
	"bitbucket.org/StephenPatrick/goevo/pop"
)

type BMDA struct {
	UMDA
	BF      []*BivariateEnv
	LastPop *pop.Population
}

func (bmda *BMDA) Adjust() Model {
	// ??????????
	// Create dependency forest
	// Generate new population from forest and frequencies
	// ??????????
	// How do you sample a 2D bivariate array
	// Do you just do it randomly
	// ??????????????
	// Combine previous population and new population
	// by direct replacement
	// Select parents
	// Calculate frequencies
	return bmda
}

func BMDAModel(opts ...Option) (Model, error) {
	var err error
	bmda := new(BMDA)
	bmda.Base, err = DefaultBase(opts...)
	// Generate initial population
	bmda.LastPop = bmda.Pop()
	bmda.LastPop.Size = bmda.learningSamples
	// Select parents?
	// I do not know why these are called parents in the pseudocode for the
	// BMDA paper, they're just the better members of the population
	subPop := bmda.selection.Select(bmda.LastPop)
	envs := make([]*env.F, len(subPop))
	for i, mem := range subPop {
		envs[i] = mem.(*UMDAIndividual).F
	}
	// Calculate univariate and bivariate frequencies
	// Univariate
	for _, e := range envs {
		bmda.F.AddF(e)
	}
	bmda.F.Divide(float64(len(subPop)))
	// At this point bmda.F holds the univariate frequencies
	// BF is the bivariate frequencies
	bmda.BF = NewFullBSBivariateEnv(envs)
	return bmda, err
}
