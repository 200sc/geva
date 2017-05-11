package eda

// CGA is an EDA that implements the Compact Genetic Algorithm
// adjustment function
type CGA struct {
	Base
}

// Adjust on a cga samples from a distribution and obtains
// the most and least fit results from that distribution,
// modifying the distribution depending on whether the most
// and least fit candidates share the same value at each index
// of their bitstring
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
	cand.Mult(cga.learningRate)
	cga.F.AddF(cand)

	cga.F.Mutate(cga.mutationRate, cga.fmutator)
	cga.learningRate = cga.lmutator(cga.learningRate)
	return cga
}

// CGAModel initializes a CGA EDA
func CGAModel(opts ...Option) (Model, error) {
	var err error
	cga := new(CGA)
	cga.Base, err = DefaultBase(opts...)
	return cga, err
}
