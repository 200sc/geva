package eda

import (
	"math"

	"bitbucket.org/StephenPatrick/goevo/env"
	"github.com/200sc/go-dist/floatrange"
)

type SHCLVND struct {
	Base
	Sigma float64
}

func BoxMueller(fr floatrange.Range) float64 {
	return math.Sqrt(-2*math.Log(fr.Poll())) * math.Sin(2*math.Pi*fr.Poll())
}

func (shc *SHCLVND) SigmaSample() *env.F {
	env := env.NewF(shc.length, 0.0)
	for i := 0; i < shc.length; i++ {
		norm := BoxMueller(floatrange.NewLinear(0, 1))
		*(*env)[i] = (norm * shc.Sigma) + *(*shc.F)[i]
	}
	return env
}

func (shc *SHCLVND) Adjust() Model {

	bcs := NewBestCandidates(shc.learningSamples)
	eCopy := shc.F.Copy()
	for i := 0; i < shc.samples; i++ {
		// We set the sample to pbil.F right now
		// as our fitness function takes in a model
		// this might change
		shc.F = shc.SigmaSample()
		bcs.Add(shc.fitness(shc), shc.F)
	}
	shc.F = eCopy
	bcsList := bcs.Slice()

	// Get average of bcsList
	mid := bcsList[0]
	for i := 1; i < len(bcsList); i++ {
		mid.AddF(bcsList[i])
	}
	mid.Divide(float64(len(bcsList)))

	mid.SubF(shc.F)
	mid.Mult(shc.learningRate)
	shc.F.AddF(mid)
	shc.Sigma = shc.lmutator(shc.Sigma)

	shc.F.Mutate(shc.mutationRate, shc.fmutator)
	return shc
}

func SHCLVNDModel(opts ...Option) (Model, error) {
	var err error
	shc := new(SHCLVND)
	shc.Base, err = DefaultBase(opts...)
	// Implicit min and max of 0 and 1
	// This .5 is taken from a paper
	shc.Sigma = (1 - 0) * 0.5
	return shc, err
}
