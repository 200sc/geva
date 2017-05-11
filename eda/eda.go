package eda

import (
	"math"

	"bitbucket.org/StephenPatrick/goevo/env"
)

// An EDA is a function which constructs an
// EDA model from a set of initialization options.
type EDA func(...Option) (Model, error)

// A Model is an iteratively adjusting EDA model
type Model interface {
	BaseModel() *Base
	Continue() bool
	Adjust() Model
	// Should be generalized
	ToEnv() *env.F
}

// A Fitness function for EDAs return an integer
// from a given model
type Fitness func(m Model) int

// Loop is the main EDA loop
func Loop(eda EDA, opts ...Option) (Model, error) {
	model, err := eda(opts...)
	if err != nil {
		return nil, err
	}
	bestFitness := math.MaxInt32
	for model.Continue() {
		model = model.Adjust()
		bm := model.BaseModel()
		if bm.trackBest {
			fitness := bm.fitness(model)
			if fitness < bestFitness {
				bm.best = bm.F.Copy()
				bestFitness = fitness
				bm.bestIteration = bm.iterations
				bm.bestFitnessEvals = bm.fitnessEvals
			}
		}
		bm.iterations++
	}
	bm := model.BaseModel()
	if bm.report != nil {
		bm.report(model)
	}
	return model, nil
}
