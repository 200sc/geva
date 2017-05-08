package eda

import "bitbucket.org/StephenPatrick/goevo/env"

// Base is a struct which all EDA models should be composed of so that
// they can use generic option functions. Future work: create several
// kinds of bases, where each base satisfies a Base interface, where
// all options will call a Set function on the interface, so that models
// that do not want as many fields as Base provides do not need to have
// wasted memory in their structs.
type Base struct {
	e            *env.F
	fitness      Fitness
	goalFitness  int
	length       int
	baseValue    float64
	randomize    bool
	learningRate float64
}

// BaseModel is a function which is used by all Options to obtain the
// base from any given model. All models must implement BaseModel.
func (b *Base) BaseModel() *Base {
	return b
}

// Option is a functional option type to be passed in variadically to model
// creation functions. All options will take the base behind a model and
// set some values on that base or otherwise manipulate the base. Options
// should be able to be called on a model in any order without the order
// changing the output model.
type Option func(Model)

// FitnessFunc is an option which sets the fitness function
func FitnessFunc(f Fitness) func(Model) {
	return func(m Model) {
		m.BaseModel().fitness = f
	}
}

// GoalFitness is an option which sets the goal fitness
func GoalFitness(gf int) func(Model) {
	return func(m Model) {
		m.BaseModel().goalFitness = gf
	}
}

// Length is an option which sets the environment length
func Length(l int) func(Model) {
	return func(m Model) {
		m.BaseModel().length = l
	}
}

// BaseValue is an option which sets the starting environment values
func BaseValue(bv float64) func(Model) {
	return func(m Model) {
		m.BaseModel().baseValue = bv
	}
}

// Randomize is an option which will randomize initial environment values
func Randomize(r bool) func(Model) {
	return func(m Model) {
		m.BaseModel().randomize = r
	}
}

// LearningRate is an option which sets the learning rate
func LearningRate(lr float64) func(Model) {
	return func(m Model) {
		m.BaseModel().learningRate = lr
	}
}
