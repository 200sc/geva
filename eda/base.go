package eda

import (
	"fmt"
	"math"
	"time"

	"github.com/200sc/geva/env"
	"github.com/200sc/geva/gevaerr"
	"github.com/200sc/geva/mut"
	"github.com/200sc/geva/pop"
	"github.com/200sc/geva/selection"

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
	name              string
	iterations        *int
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
	bestFitnessEvals  *int
	bestIteration     *int
	bestFitness       *int
	attemptsAfterBest int
	best              *env.F
	startTime         time.Time
	// Ensure booleans are kept at the end (or in groups of 8)
	// for memory efficiency
	randomize    bool
	trackBest    bool
	trackFitness bool
	trackTime    bool
}

// DefaultBase initializes some base fields to non-automatic zero values
func DefaultBase(opts ...Option) (Base, error) {
	b := new(Base)
	b.fmutator = mut.None()
	b.lmutator = mut.None()
	b.length = 1
	b.iterations = new(int)
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
		return *b, gevaerr.InvalidParamError{"length"}
	}
	if b.samples <= 0 {
		return *b, gevaerr.InvalidParamError{"samples"}
	}
	if b.learningSamples <= 0 {
		return *b, gevaerr.InvalidParamError{"learningSamples"}
	}
	if !b.valueRange.InRange(b.baseValue) {
		return *b, gevaerr.InvalidParamError{"baseValue"}
	}
	if b.trackFitness {
		ofitness := b.fitness
		fevals := b.fitnessEvals
		b.fitness = func(e *env.F) int {
			(*fevals)++
			return ofitness(e)
		}
	}
	if b.trackBest {
		b.bestFitness = new(int)
		b.bestFitnessEvals = new(int)
		b.bestIteration = new(int)
		b.best = env.NewF(1, 0.0)
		*b.bestFitness = math.MaxInt32
		ofitness := b.fitness
		b.fitness = func(e *env.F) int {
			// These are pointers because the function closure that is
			// generated doesn't keep track of b, it looks like an
			// optimization method means that when you write b.best
			// the referred to value is not (look through b and find best)
			// but a environment local version of b.best. An alternative
			// approach would be to write these functions in a way that they
			// obtain b at the beginning to force re-looking for b.best, etc
			f := ofitness(e)
			if f < *b.bestFitness {
				*b.best = *e
				*b.bestFitness = f
				*b.bestIteration = *b.iterations
				if b.fitnessEvals != nil {
					*b.bestFitnessEvals = *b.fitnessEvals
				}
			}
			return f
		}
	}
	b.F = env.NewF(b.length, b.baseValue)
	if b.randomize {
		b.F.RandomizeSingle(0.0, 1.0)
	}
	if b.trackTime {
		b.startTime = time.Now()
	}
	return *b, nil
}

// BaseModel is a function which is used by all Options to obtain the
// base from any given model. All models must implement BaseModel.
func (b *Base) BaseModel() *Base {
	return b
}

// Fitness is shorthand for fitness(b.F)
func (b *Base) Fitness() int {
	return b.fitness(b.ToEnv())
}

// DefReport is the Default Report function
func DefReport(m Model) {
	bm := m.BaseModel()
	fmt.Println("Test", bm.name)
	fmt.Println("Iterations taken:", *bm.iterations)
	if bm.best != nil {
		fmt.Println("Best Model:", bm.best)
		fmt.Println("Best Fitness:", *bm.bestFitness)
		fmt.Println("Iteration of best model:", *bm.bestIteration)
	}
	if bm.fitnessEvals != nil {
		fmt.Println("Fitness Evaluations:", *bm.fitnessEvals)
		fmt.Println("Fitness Evals at best model iteration:", *bm.bestFitnessEvals)
	}
	if bm.trackTime {
		fmt.Println("Time taken:", time.Since(bm.startTime))
	}
}

// DefContinue is the Default Continue function
func DefContinue(m Model) bool {
	b := m.BaseModel()
	fitness := b.Fitness()
	//fmt.Println("Iteration", *b.iterations, "Fitness:", fitness)
	return fitness > b.goalFitness && *b.iterations < b.maxIterations
}

// Continue is a wrapper around b.cont
func (b *Base) Continue() bool {
	return b.cont(b)
}

// Adjust is the default Adjust method, which is expected to be overwritten
func (b *Base) Adjust() Model {
	return b
}

// Mutate is the default mutation function
func (b *Base) Mutate() {
	b.F.Mutate(b.mutationRate, b.fmutator)
	b.learningRate = b.lmutator(b.learningRate)
}

// GenIndices is a utility function to generate a list of ints 0 ... b.length
func (b *Base) GenIndices() []int {
	available := make([]int, b.length)
	for i := range available {
		available[i] = i
	}
	return available
}

// Pop generates a population of EnvInds from sampling base.F
func (b *Base) Pop() *pop.Population {

	// Generate a population of size samples by sampling umda
	p := new(pop.Population)
	p.Members = make([]pop.Individual, b.samples)
	for i := 0; i < b.samples; i++ {
		p.Members[i] = NewEnvInd(b.F)
	}

	// Generate fitnesses for the population
	p.Fitnesses = make([]int, b.samples)
	for i := 0; i < b.samples; i++ {
		p.Fitnesses[i] = b.fitness(p.Members[i].(*EnvInd).F)
	}

	p.Size = b.samples
	return p
}

// SelectLearning selects b.learningSamples members from the given population
func (b *Base) SelectLearning(p *pop.Population) []pop.Individual {
	oldSize := p.Size
	p.Size = b.learningSamples
	selected := b.selection.Select(p)
	p.Size = oldSize
	return selected
}
