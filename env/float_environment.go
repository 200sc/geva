package env

import (
	"math"
	"math/rand"

	"bitbucket.org/StephenPatrick/goevo/mut"

	"github.com/200sc/go-dist/floatrange"
)

// F returns a float-valued environment
type F []*float64

// ToEnv is an assistant function for structs composed with
// F for interfaces to obtain their F
func (env *F) ToEnv() *F {
	return env
}

// Copy returns a copy of F.
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

// New creates a combined environment by the given envDiff inputs
// and returns a new environment pointer
func (env *F) New(envDiff []float64) *F {
	newEnv := make(F, len(*env))
	for i, f := range envDiff {
		newEnv[i] = new(float64)
		*newEnv[i] = *(*env)[i] + f
	}
	return &newEnv
}

// Diff returns the absolute difference between the given environment
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

// DiffSingle returns the sum of the difference between f and each element in
// the environment
func (env *F) DiffSingle(f float64) int {
	diff := 0.0
	for _, f2 := range *env {
		if *f2 != 0.0 {
			diff += math.Abs(*f2 - f)
		}
	}
	return int(diff)
}

// SetAll sets each value behind the environment to f
func (env *F) SetAll(f float64) {
	for i := range *env {
		*(*env)[i] = f
	}
}

// Randomize sets env[i] to a unfiorm sample between min[i] and max[i] for all i
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

// RandomizeSingle sets env[i] to a uniform sample between mn and max for all i
func (env *F) RandomizeSingle(mn float64, mx float64) {
	diff := mx - mn
	for i := 0; i < len(*env); i++ {
		*(*env)[i] = (diff * rand.Float64()) + mn
	}
}

// ToIntRandom treats the environment as a set of probabilities, and creates
// an env.I where each element has a env[i] percent chance of being 1, else 0.
func (env *F) ToIntRandom() *I {
	envI := NewI(len(*env), 0)
	for i, f := range *env {
		if rand.Float64() < *f {
			*(*envI)[i] = 1
		}
	}
	return envI
}

// Reinforce updates each value in env to be learningRate closer to env2.
func (env *F) Reinforce(env2 *F, learningRate float64) {
	learningRate = floatrange.NewLinear(0, 1).EnforceRange(learningRate)
	negRate := 1.0 - learningRate
	for i, f := range *env {
		*f = (*f * negRate) + (*(*env2)[i] * learningRate)
	}
}

func (env *F) Mutate(chance float64, mutator mut.FloatMutator) {
	for _, f := range *env {
		if rand.Float64() < chance {
			*f = mutator(*f)
		}
	}
}

// NewF creates an environment
func NewF(size int, baseVal float64) *F {
	env := make(F, size)
	for i := 0; i < size; i++ {
		env[i] = new(float64)
		*env[i] = baseVal
	}
	return &env
}
