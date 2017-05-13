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
	fitnessEvals      *int
	bestFitnessEvals  int
	bestIteration     int
	attemptsAfterBest int
	best              *env.F
	// Ensure booleans are kept at the end (or in groups of 8)
	// for memory efficiency
	randomize    bool
	trackBest    bool
	trackFitness bool
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
	b.valueRange = floatrange.NewLinear(0, 1)
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
	if b.trackFitness {
		ofitness := b.fitness
		fevals := b.fitnessEvals
		b.fitness = func(e *env.F) int {
			(*fevals)++
			return ofitness(e)
		}
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

func (b *Base) Fitness() int {
	return b.fitness(b.ToEnv())
}

// DefReport is the Default Report function
func DefReport(m Model) {
	bm := m.BaseModel()
	fmt.Println("Iterations taken:", bm.iterations)
	if bm.best != nil {
		fmt.Println("Best Model:", bm.best)
		fmt.Println("Best Fitness:", bm.fitness(bm.best))
		fmt.Println("Iteration of best model:", bm.bestIteration)
	}
	if bm.fitnessEvals != nil {
		fmt.Println("Fitness Evaluations:", *bm.fitnessEvals)
		fmt.Println("Fitness Evals at best model iteration:", bm.bestFitnessEvals)
	}
}

// DefContinue is the Default Continue function
func DefContinue(m Model) bool {
	b := m.BaseModel()
	fitness := b.Fitness()
	fmt.Println("Iteration", b.iterations, "Fitness:", fitness)
	return fitness > b.goalFitness && b.iterations < b.maxIterations
}

// Continue is a wrapper around b.cont
func (b *Base) Continue() bool {
	return b.cont(b)
}

// Adjust is the default Adjust method, which is expected to be overwritten
func (b *Base) Adjust() Model {
	return b
}

// GenIndices is a utility function to generate a list of ints 0 ... b.length
func (b *Base) GenIndices() []int {
	available := make([]int, b.length)
	for i := range available {
		available[i] = i
	}
	return available
}
