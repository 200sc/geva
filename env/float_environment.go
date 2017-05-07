package env

import (
	"math"
	"math/rand"
)

type F []*float64

func (env *F) Copy() *F {
	newEnv := make(F, len(*env))
	for i, f := range *env {
		var v float64
		if f == nil {
			v = 0.0
		} else {
			v = *f
		}
		newEnv[i] = new(float64)
		*newEnv[i] = v
	}
	return &newEnv
}

// Creates a combined environment by the given envDiff inputs
// and returns a new environment pointer
func (env *F) New(envDiff []float64) *F {
	newEnv := make(F, len(*env))
	for i, f := range envDiff {
		newEnv[i] = new(float64)
		*newEnv[i] = *(*env)[i] + f
	}
	return &newEnv
}

// Returns the absolute difference between the given environment
// and the passed in expectations. 0 in envDiff is treated as
// insignificant.
func (env *F) Diff(envDiff []float64) int {
	diff := 0.0
	for i, f := range envDiff {
		if f != 0.0 {
			diff += math.Abs(*(*env)[i] - f)
		}
	}
	return int(diff)
}

// DiffSingle returns the difference between
func (env *F) DiffSingle(check float64) int {
	diff := 0.0
	for i, f := range *env {
		if *f != 0.0 {
			diff += math.Abs(*(*env)[i] - check)
		}
	}
	return int(diff)
}

func (env *F) SetAll(f float64) {
	for i := range *env {
		*(*env)[i] = f
	}
}

func (env *F) Randomize(min *F, max *F) {
	// Probably shouldn't do this
	// if len(env) != len(min) || len(env) != len(max) {
	// 	panic("It won't do what you want")
	// }
	for i := 0; i < len(*env); i++ {
		mn := *(*min)[i]
		mx := *(*max)[i]
		diff := mx - mn
		*(*env)[i] = (diff * rand.Float64()) + mn
	}
}

func (env *F) RandomizeSingle(mn float64, mx float64) {
	diff := mx - mn
	for i := 0; i < len(*env); i++ {
		*(*env)[i] = (diff * rand.Float64()) + mn
	}
}

func NewF(size int, baseVal float64) *F {
	env := make(F, size)
	for i := 0; i < size; i++ {
		env[i] = new(float64)
		*env[i] = baseVal
	}
	return &env
}
