package eda

import "bitbucket.org/StephenPatrick/goevo/env"

// An EDA is a function which constructs an
// EDA model from a set of initialization options.
type EDA func(...Option) (Model, error)

// A Model is an iteratively adjusting EDA model
type Model interface {
	BaseModel() *Base
	Continue() bool
	Adjust(samples int) Model
	// Should be generalized
	ToEnv() *env.F
}

// A Fitness function for EDAs return an integer
// from a given model
type Fitness func(m Model) int

// Loop is the main EDA loop
func Loop(eda EDA, samples int, opts ...Option) (Model, error) {
	model, err := eda(opts...)
	if err != nil {
		return nil, err
	}
	for model.Continue() {
		model = model.Adjust(samples)
	}
	return model, nil
}
