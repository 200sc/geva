package gp

import (
	"math"
)

// The index of each environment element is insignificant.
// GPs are expected to learn which indexes are useful for
// each operation.
// The int that each element refers to means something
// particular for each problem. GPs are also expected to
// learn what, of these, is important.
type Environment []*int

func (env Environment) Copy() *Environment {
	newEnv := make(Environment, len(env))
	for i, f := range env {
		var v int
		if f == nil {
			v = 0
		} else {
			v = *f
		}
		newEnv[i] = new(int)
		*newEnv[i] = v
	}
	return &newEnv
}

// Creates a combined environment by the given envDiff inputs
// and returns a new environment pointer
func (env Environment) New(envDiff []float64) *Environment {
	newEnv := make(Environment, len(env))
	for i, f := range envDiff {
		newEnv[i] = new(int)
		*newEnv[i] = int(math.Ceil(float64(*env[i]) + f))
	}
	return &newEnv
}

// Returns the absolute difference between the given environment
// and the passed in expectations. 0 in envDiff is treated as
// insinificant.
func (env Environment) Diff(envDiff []float64) (diff int) {
	for i, f := range envDiff {
		if f != 0.0 {
			diff += int(math.Ceil(float64(*env[i]) - f))
		}
	}
	return
}

func NewEnvironment(size int, baseVal int) Environment {
	env := make(Environment, size)
	for i := 0; i < size; i++ {
		env[i] = new(int)
		*env[i] = baseVal
	}
	return env
}
