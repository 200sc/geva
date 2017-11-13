package eda

import "github.com/200sc/geva/env"

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
		for j, f := range *(ind.(*EnvInd).F) {
			newenv.Set(j, newenv.Get(j)+*f)
		}
	}
	// divide the resulting sums by the total number of samples
	// to obtain a new umda
	newenv.Divide(float64(len(subPop)))
	umda.F = newenv

	return umda
}

// UMDAModel initializes a UMDA EDA
func UMDAModel(opts ...Option) (Model, error) {
	var err error
	umda := new(UMDA)
	umda.Base, err = DefaultBase(opts...)
	return umda, err
}
