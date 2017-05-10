package eda

import (
	"fmt"
	"time"

	"bitbucket.org/StephenPatrick/goevo/env"
	"bitbucket.org/StephenPatrick/goevo/evoerr"
)

type CGA struct {
	Base
}

func (cga *CGA) Continue() bool {
	fitness := cga.fitness(cga)
	fmt.Println(fitness, cga.goalFitness)
	time.Sleep(1 * time.Second)
	return fitness > cga.goalFitness
}

func (cga *CGA) Adjust() Model {

	bcs := NewBestCandidates(2)
	eCopy := cga.F.Copy()
	for i := 0; i < 2; i++ {
		// We set the sample to cga.F right now
		// as our fitness function takes in a model
		// this might change
		cga.F = GetSample(eCopy)
		bcs.Add(cga.fitness(cga), cga.F)
	}

	cga.F = eCopy
	bcsList := bcs.Slice()
	cand := bcsList[0]
	
	cand.SubF(bcsList[1])
	fmt.Println(cand)
	cand.Mult(cga.learningRate)
	cga.F.AddF(cand)

	cga.F.Mutate(cga.mutationRate, cga.fmutator)
	cga.learningRate = cga.lmutator(cga.learningRate)
	return cga
}

func CGAModel(opts ...Option) (Model, error) {
	cga := new(CGA)
	cga.Base = DefaultBase()
	for _, opt := range opts {
		opt(cga)
	}
	if cga.length <= 0 {
		return nil, evoerr.InvalidLengthError{}
	}
	cga.F = env.NewF(cga.length, cga.baseValue)
	if cga.randomize {
		cga.F.RandomizeSingle(0.0, 1.0)
	}
	return cga, nil
}
