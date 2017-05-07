package eda

import "bitbucket.org/StephenPatrick/goevo/env"

type PBIL struct {
	Base
}

func (pbil *PBIL) ShouldContinue() bool {
	fitness := pbil.fitness(pbil)
	return fitness < pbil.goalFitness
}

func (pbil *PBIL) Adjust(samples int) Model {
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
}

type InvalidLengthError struct{}

func (ile InvalidLengthError) Error() string {
	return "The length given was less than or equal to zero"
}
