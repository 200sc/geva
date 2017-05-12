package eda

import (
	"math"
	"math/rand"

	"bitbucket.org/StephenPatrick/goevo/env"
	"bitbucket.org/StephenPatrick/goevo/pop"
)

// BMDA represents the Bivariate Marginal Distribution Algorithm
type BMDA struct {
	UMDA
	BF       FullBivariateEnv
	LastPop  *pop.Population
	Roots    []int
	Children [][]int
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
	// Reset the previous forest
	bmda.Roots = []int{}
	bmda.Children = make([][]int, bmda.length)
	for i := 0; i < bmda.length; i++ {
		bmda.Children[i] = []int{}
	}
	available := bmda.GenIndices()
	used := []int{}
	for {
		// choose a random index
		i := rand.Intn(len(available))
		chosen := available[i]
		// create a new tree in the forest starting at chosen
		bmda.Roots = append(bmda.Roots, chosen)
		for len(available) > 1 {
			// remove the chosen index from available
			available = append(available[:i], available[i+1:]...)
			used = append(used, chosen)

			// Find the most dependant pairing given our set of used
			// and available indices
			maxChi2 := 0.0
			var parent int
			for _, v := range available {
				for _, v2 := range used {
					chi2 := bmda.ChiSquared(v, v2)
					if chi2 > maxChi2 {
						maxChi2 = chi2
						chosen = v
						parent = v2
					}
				}
			}

			// If no dependencies exist break to the outer loop.
			if maxChi2 < 3.84 {
				break
			}
			bmda.Children[parent] = append(bmda.Children[parent], chosen)
		}
		if len(available) == 0 {
			break
		}
	}
	// Generate new population from forest and frequencies

	// How do you sample a 2D bivariate array
	// No you use the foest we just spent way too long making
	// So we pick random (all) roots and go down through all of their children
	// using their parents as the elements they are dependant on

	// Combine previous population and new population by direct replacement
	// of I guess arbitrary elements
	bmda.UpdateFromPop()
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
	bmda.UpdateFromPop()
	return bmda, err
}

func (bmda *BMDA) UpdateFromPop() {
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
}
