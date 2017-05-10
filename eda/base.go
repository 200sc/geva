package eda

import (
	"fmt"

	"bitbucket.org/StephenPatrick/goevo/env"
	"bitbucket.org/StephenPatrick/goevo/evoerr"
	"bitbucket.org/StephenPatrick/goevo/mut"
	"bitbucket.org/StephenPatrick/goevo/selection"
)

// Base is a struct which all EDA models should be composed of so that
// they can use generic option functions. Future work: create several
// kinds of bases, where each base satisfies a Base interface, where
// all options will call a Set function on the interface, so that models
// that do not want as many fields as Base provides do not need to have
// wasted memory in their structs.
type Base struct {
	*env.F
	fitness         Fitness
	goalFitness     int
	length          int
	baseValue       float64
	learningRate    float64
	mutationRate    float64
	lmutator        mut.FloatMutator
	fmutator        mut.FloatMutator
	samples         int
	learningSamples int
	selection       selection.Method
	randomize       bool
}

// DefaultBase initializes some base fields to non-automatic zero values
func DefaultBase(opts ...Option) (Base, error) {
	b := new(Base)
	b.fmutator = mut.None()
	b.lmutator = mut.None()
	b.length = 1
	b.samples = 1
	b.learningSamples = 1
	for _, opt := range opts {
		opt(b)
	}
	if b.length <= 0 {
		return *b, evoerr.InvalidLengthError{}
	}
	b.F = env.NewF(b.length, b.baseValue)
	if b.randomize {
		b.F.RandomizeSingle(0.0, 1.0)
	}
	return *b, nil
}

// BaseModel is a function which is used by all Options to obtain the
// base from any given model. All models must implement BaseModel.
func (b *Base) BaseModel() *Base {
	return b
}

// Default Continue
func (b *Base) Continue() bool {
	fitness := b.fitness(b)
	fmt.Println(fitness, b.goalFitness)
	return fitness > b.goalFitness
}

// Default Adjust
func (b *Base) Adjust() Model {
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

func Samples(s int) func(Model) {
	return func(m Model) {
		m.BaseModel().samples = s
	}
}

func LearningSamples(l int) func(Model) {
	return func(m Model) {
		m.BaseModel().learningSamples = l
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

// MutationRate is an option which sets the mutation rate
func MutationRate(mr float64) func(Model) {
	return func(m Model) {
		m.BaseModel().mutationRate = mr
	}
}

// FMutator sets a model's float mutator
func FMutator(mtr mut.FloatMutator) func(Model) {
	return func(m Model) {
		m.BaseModel().fmutator = mtr
	}
}

// LMutator sets a model's learning rate mutator
func LMutator(mtr mut.FloatMutator) func(Model) {
	return func(m Model) {
		m.BaseModel().lmutator = mtr
	}
}

func SelectionMethod(sm selection.Method) func(Model) {
	return func(m Model) {
		m.BaseModel().selection = sm
	}
}
