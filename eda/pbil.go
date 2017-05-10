package eda

import (
	"math/rand"

	"bitbucket.org/StephenPatrick/goevo/env"
	"bitbucket.org/StephenPatrick/goevo/evoerr"
)

type PBIL struct {
	Base
}

func (pbil *PBIL) Continue() bool {
	fitness := pbil.fitness(pbil)
	//fmt.Println(fitness, pbil.goalFitness)
	return fitness > pbil.goalFitness
}

func GetSample(e *env.F) *env.F {
	sample := e.Copy()
	for _, f := range *sample {
		if rand.Float64() <= *f {
			*f = 1
		} else {
			*f = 0
		}
	}
	return sample
}

func (pbil *PBIL) Adjust() Model {

	bcs := NewBestCandidates(pbil.learningSamples)
	eCopy := pbil.F.Copy()
	for i := 0; i < pbil.samples; i++ {
		// We set the sample to pbil.F right now
		// as our fitness function takes in a model
		// this might change
		pbil.F = GetSample(eCopy)
		bcs.Add(pbil.fitness(pbil), pbil.F)
	}
	// Also could add a worst candidate and a negative learning rate
	pbil.F = eCopy
	bcsList := bcs.Slice()
	// Hypothetically bcsList has a length equal to
	// pbil.learningSamples but if samples < learningSamples
	// this this case ensures we still learn a total of learningRate.
	//fmt.Println("Learning rate:", pbil.learningRate)
	realRate := pbil.learningRate / float64(len(bcsList))
	for _, cand := range bcsList {
		pbil.F.Reinforce(cand, realRate)
	}
	pbil.F.Mutate(pbil.mutationRate, pbil.fmutator)
	pbil.learningRate = pbil.lmutator(pbil.learningRate)
	return pbil
}

func PBILModel(opts ...Option) (Model, error) {
	pbil := new(PBIL)
	pbil.Base = DefaultBase()
	for _, opt := range opts {
		opt(pbil)
	}
	if pbil.length <= 0 {
		return nil, evoerr.InvalidLengthError{}
	}
	pbil.F = env.NewF(pbil.length, pbil.baseValue)
	if pbil.randomize {
		pbil.F.RandomizeSingle(0.0, 1.0)
	}
	return pbil, nil
}
