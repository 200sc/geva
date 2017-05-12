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
	bcs := NewBestCandidates(cga, cga.samples, nil)
	bcs.Front.F.SubF(bcs.Back.F).Mult(cga.learningRate)
	cga.F.AddF(bcs.Front.F)
	return cga
}

// CGAModel initializes a CGA EDA
func CGAModel(opts ...Option) (Model, error) {
	var err error
	cga := new(CGA)
	cga.Base, err = DefaultBase(opts...)
	return cga, err
}
