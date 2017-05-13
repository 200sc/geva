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
