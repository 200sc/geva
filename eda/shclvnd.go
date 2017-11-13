package eda

import (
	"github.com/200sc/geva/eda/stat"
	"github.com/200sc/geva/env"
	"github.com/200sc/go-dist/floatrange"
)

// SHCLVND represents the Stochastic Hill Climbing with Learning by Vectors
// of Normal Distributions algorithm.
type SHCLVND struct {
	Base
	Sigma float64
}

// SigmaSample is a way to sample an environment based on a standard
// deviation spread from an average
func (shc *SHCLVND) SigmaSample() *env.F {
	env := env.NewF(shc.length, 0.0)
	for i := 0; i < shc.length; i++ {
		norm := stat.BoxMueller(floatrange.NewLinear(0, 1))
		env.Set(i, shc.F.Get(i)+(norm*shc.Sigma))
	}
	return env
}

// Adjust on a SHCLVND polls learningSamples best samples and creates an
// average candidate from the collected samples. Then the distribution of
// the SHCLVND is reinforced closer to the average candidate, and the standard
// deviation used to generate the samples is reduced.
func (shc *SHCLVND) Adjust() Model {
	bcs := NewBestCandidates(shc, shc.learningSamples, shc.SigmaSample)
	mid := env.AverageF(bcs.Slice()...)

	mid.SubF(shc.F).Mult(shc.learningRate)
	shc.F.AddF(mid)
	shc.Sigma = shc.lmutator(shc.Sigma)
	return shc
}

// SHCLVNDModel initializes a SHCLVND EDA
func SHCLVNDModel(opts ...Option) (Model, error) {
	var err error
	shc := new(SHCLVND)
	shc.Base, err = DefaultBase(opts...)
	// Implicit min and max of 0 and 1
	// This .5 is taken from the paper, "Stochastic hill climbing with learning
	// by vectors of normal distributions"
	shc.Sigma = (1 - 0) * 0.5
	return shc, err
}
