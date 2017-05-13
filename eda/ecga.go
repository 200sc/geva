package eda

type ECGA struct {
	Base
	Pop
	MPM
}

func (ecga *ECGA) Adjust() Model {
	// Crossover with elites
	// Selection
	// MPM Model using MDL
	return cga
}

func ECGAModel(opts ...Option) (Model, error) {
	var err error
	ecga := new(ECGA)
	ecga.Base, err = DefaultBase(opts...)
	// Random Pop
	// Selection
	// Build MPM Model using MDL???
	return ecga, err
}

// ECGA isn't terribly hard to do-- we already have crossover and selection
// covered, and elitism is fine but we should probably give the other models
// who use populations elitism if we give ECGA elitism. The one hard thing
// is the MPM model which in itself isn't hard but requires that we split
// the problem space up into building blocks, and requires that we learn these
// building blocks, I guess, although I'm not certain on that
//
// Yes, according to paper 1 we begin the MDM model with the assumption that
// every index is its own building block, and sees if merging all pairs is 
// helpful. 
//
// Papers : https://pdfs.semanticscholar.org/eeee/a9fdade929cb3fc9a99631d3541ef7005079.pdf
// http://www.kumarasastry.com/wp-content/files/2000026.pdf
