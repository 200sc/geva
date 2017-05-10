package eda

import "fmt"

type CGA struct {
	Base
}

func (cga *CGA) Adjust() Model {

	bcs := NewBestCandidates(cga.samples)
	eCopy := cga.F.Copy()
	for i := 0; i < cga.samples; i++ {
		// We set the sample to cga.F right now
		// as our fitness function takes in a model
		// this might change
		cga.F = GetSample(eCopy)
		bcs.Add(cga.fitness(cga), cga.F)
	}

	cga.F = eCopy
	cand := bcs.Front.F

	cand.SubF(bcs.Back.F)
	fmt.Println(cand)
	cand.Mult(cga.learningRate)
	cga.F.AddF(cand)

	cga.F.Mutate(cga.mutationRate, cga.fmutator)
	cga.learningRate = cga.lmutator(cga.learningRate)
	return cga
}

func CGAModel(opts ...Option) (Model, error) {
	var err error
	cga := new(CGA)
	cga.Base, err = DefaultBase(opts...)
	return cga, err
}
