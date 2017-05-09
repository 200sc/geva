package eda

import (
	"math"
	"math/rand"

	"bitbucket.org/StephenPatrick/goevo/env"
	"bitbucket.org/StephenPatrick/goevo/evoerr"
)

type PBIL struct {
	Base
}

func (pbil *PBIL) ShouldContinue() bool {
	fitness := pbil.fitness(pbil)
	return fitness < pbil.goalFitness
}

func (pbil *PBIL) Adjust(samples int) Model {
	bestCandidateFitness := math.MaxInt32
	var bestCandidate *env.F
	eCopy := pbil.F.Copy()
	for i := 0; i < samples; i++ {
		pbil.F = pbil.F.Copy()
		for j, f := range *pbil.F {
			if rand.Float64() < *f {
				*(*pbil.F)[j] = 1
			} else {
				*(*pbil.F)[j] = 0
			}
		}
		f := pbil.fitness(pbil)
		if f < bestCandidateFitness {
			bestCandidateFitness = f
			bestCandidate = pbil.F.Copy()
		}
	}
	pbil.F = eCopy
	pbil.F.Reinforce(bestCandidate, pbil.learningRate)
	pbil.F.Mutate(pbil.mutationRate, pbil.fmutator)
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
