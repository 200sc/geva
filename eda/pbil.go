package eda

import (
	"math"
	"math/rand"

	"bitbucket.org/StephenPatrick/goevo/env"
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
	eCopy := pbil.e.Copy()
	for i := 0; i < samples; i++ {
		pbil.e = pbil.e.Copy()
		for j, f := range *pbil.e {
			if rand.Float64() < *f {
				*(*pbil.e)[j] = 1
			} else {
				*(*pbil.e)[j] = 0
			}
		}
		f := pbil.fitness(pbil)
		if f < bestCandidateFitness {
			bestCandidateFitness = f
			bestCandidate = pbil.e.Copy()
		}
	}
	pbil.e = eCopy
	pbil.e.Learn(eCopy, pbil.learningRate)
	//pbil.e.Mutate()
	return pbil
}

func (pbil *PBIL) ToEnv() *env.F {
	return pbil.e
}

func PBILModel(opts ...Option) (Model, error) {
	pbil := new(PBIL)
	for _, opt := range opts {
		opt(pbil)
	}
	if pbil.length <= 0 {
		return nil, InvalidLengthError{}
	}
	pbil.e = env.NewF(pbil.length, pbil.baseValue)
	if pbil.randomize {
		pbil.e.RandomizeSingle(0.0, 1.0)
	}
	return pbil, nil
}

type InvalidLengthError struct{}

func (ile InvalidLengthError) Error() string {
	return "The length given was less than or equal to zero"
}
