package eda

import (
	"fmt"
	"math"

	"bitbucket.org/StephenPatrick/goevo/env"
	"bitbucket.org/StephenPatrick/goevo/evoerr"
	"bitbucket.org/StephenPatrick/goevo/mut"
	"bitbucket.org/StephenPatrick/goevo/selection"

	"github.com/200sc/go-dist/floatrange"
)

// Base is a struct which all EDA models should be composed of so that
// they can use generic option functions. Future work: create several
// kinds of bases, where each base satisfies a Base interface, where
// all options will call a Set function on the interface, so that models
// that do not want as many fields as Base provides do not need to have
// wasted memory in their structs.
type Base struct {
	*env.F
	iterations        int
	maxIterations     int
	fitness           Fitness
	goalFitness       int
	length            int
	valueRange        floatrange.Range
	baseValue         float64
	learningRate      float64
	mutationRate      float64
	lmutator          mut.FloatMutator
	fmutator          mut.FloatMutator
	samples           int
	learningSamples   int
	selection         selection.Method
	cont              func(Model) bool
	report            func(Model)
	bestIteration     int
	attemptsAfterBest int
	best              *env.F
	randomize         bool
	trackBest         bool
}

// DefaultBase initializes some base fields to non-automatic zero values
func DefaultBase(opts ...Option) (Base, error) {
	b := new(Base)
	b.fmutator = mut.None()
	b.lmutator = mut.None()
	b.length = 1
	b.samples = 1
	b.learningSamples = 1
	b.maxIterations = math.MaxInt32
	b.valueRange = floatrange.NewLinear(0, math.MaxFloat64)
	b.cont = DefContinue
	b.report = DefReport
	for _, opt := range opts {
		opt(b)
	}
	if b.length <= 0 {
		return *b, evoerr.InvalidParamError{"length"}
	}
	if b.samples <= 0 {
		return *b, evoerr.InvalidParamError{"samples"}
	}
	if b.learningSamples <= 0 {
		return *b, evoerr.InvalidParamError{"learningSamples"}
	}
	if !b.valueRange.InRange(b.baseValue) {
		return *b, evoerr.InvalidParamError{"baseValue"}
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

// DefReport is the Default Report function
func DefReport(m Model) {
	bm := m.BaseModel()
	fmt.Println("Iterations taken:", bm.iterations)
	if bm.best != nil {
		fmt.Println("Best Model:", bm.best)
		bm.F = bm.best
		fmt.Println("Best Fitness:", bm.fitness(bm))
		fmt.Println("Iteration of best model:", bm.bestIteration)
	}
}

// DefContinue is the Default Continue function
func DefContinue(m Model) bool {
	b := m.BaseModel()
	fitness := b.fitness(b)
	//fmt.Println(fitness, b.goalFitness)
	return fitness > b.goalFitness && b.iterations < b.maxIterations
}

func (b *Base) Continue() bool {
	return b.cont(b)
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

// Samples sets the number of samples the algorithm will use.
func Samples(s int) func(Model) {
	return func(m Model) {
		m.BaseModel().samples = s
	}
}

// LearningSamples sets the number of learning samples the algorithm will use.
// learning samples should be less than samples.
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

// SelectionMethod sets the algorithm's population selection method for
// algorithms that use selection
func SelectionMethod(sm selection.Method) func(Model) {
	return func(m Model) {
		m.BaseModel().selection = sm
	}
}

// StopCondition replaces the default stop condition (goal fitness reached
// or maximum iteration reached)
func StopCondition(c func(m Model) bool) func(Model) {
	return func(m Model) {
		m.BaseModel().cont = c
	}
}

// TrackBest sets whether or not the algorithm should keep track of
// the best model it has found so far
func TrackBest(b bool) func(Model) {
	return func(m Model) {
		m.BaseModel().trackBest = b
	}
}

// AttemptsAfterBest terminates an EDA if it fails to improve on its best
// model after i iterations
func AttemptsAfterBest(i int) func(Model) {
	return func(m Model) {
		bm := m.BaseModel()
		bm.trackBest = true
		bm.attemptsAfterBest = i
		oldCont := bm.cont
		bm.cont = func(m Model) bool {
			bm := m.BaseModel()
			return oldCont(m) && (bm.bestIteration+bm.attemptsAfterBest) > bm.iterations
		}
	}
}

// MaxIterations sets the maximum number of iterations the algorithm
// will evolve before preemptively stopping
func MaxIterations(j int) func(Model) {
	return func(m Model) {
		m.BaseModel().maxIterations = j
	}
}

// And can be used to chain together multiple settings in order to
// prepackage common settings among several algorithms.
// Consider: should this be called Cons?
func And(a, b Option) func(Model) {
	return func(m Model) {
		a(m)
		b(m)
	}
}
