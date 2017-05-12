package eda

import (
	"bitbucket.org/StephenPatrick/goevo/mut"
	"bitbucket.org/StephenPatrick/goevo/selection"
)

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

// TrackFitnessRuns tells the model to track how many times it calls its
// fitness function.
func TrackFitnessRuns(b bool) func(Model) {
	return func(m Model) {
		m.BaseModel().trackFitness = b
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
func And(opts ...Option) func(Model) {
	return func(m Model) {
		for _, opt := range opts {
			opt(m)
		}
	}
}
