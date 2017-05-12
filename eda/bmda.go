package eda

import (
	"math"

	"bitbucket.org/StephenPatrick/goevo/env"
	"bitbucket.org/StephenPatrick/goevo/pop"
)

// BMDA represents the Bivariate Marginal Distribution Algorithm
type BMDA struct {
	UMDA
	BF      FullBivariateEnv
	LastPop *pop.Population
}

// ChiSquared on a bmda calculates a chi^2 value
// where the observed values are the unconditional bivariate
// probabilities and the expected values are the univariate
// probabilities multipled together.
func (bmda *BMDA) ChiSquared(a, b int) float64 {

	// p(a=t)
	at := bmda.F.Get(a)
	bt := bmda.F.Get(b)
	af := 1 - at
	bf := 1 - bt

	// p(a=t|b=t)
	catbt := bmda.BF[a].bf[0].Get(b)
	catbf := bmda.BF[a].bf[1].Get(b)
	cafbt := 1 - catbt
	cafbf := 1 - catbf

	// p(a=t|b=t) = p(a=t,b=t) / p(b=t) so
	// p(a=t,b=t)
	atbt := catbt * bt
	atbf := catbf * bf
	afbt := cafbt * bt
	afbf := cafbf * bf

	chi2 := 0.0
	fsamp := float64(bmda.samples)

	xptd1 := fsamp * at * bt
	chi2 += math.Pow((fsamp*atbt)-xptd1, 2) / xptd1

	xptd2 := fsamp * at * bf
	chi2 += math.Pow((fsamp*atbf)-xptd2, 2) / xptd2

	xptd3 := fsamp * af * bt
	chi2 += math.Pow((fsamp*afbt)-xptd3, 2) / xptd3

	xptd4 := fsamp * af * bf
	chi2 += math.Pow((fsamp*afbf)-xptd4, 2) / xptd4

	return chi2
}

// Dependant uses bmda.ChiSquared to determine whether a is dependant on b.
// This uses a constant taken from the paper "The Bivariate Marginal Distribution
// Algorithm".
func (bmda *BMDA) Dependant(a, b int) bool {
	return bmda.ChiSquared(a, b) > 3.84
}

// Adjust on a BMDA is incomplete
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

// BMDAModel returns an initialized BMDA
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
	// BF is the conditional bivariate frequencies
	bmda.BF = NewFullBSBivariateEnv(envs)
	return bmda, err
}
