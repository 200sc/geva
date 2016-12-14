package env

import (
	//"fmt"
	"math"
)

// The index of each environment element is insignificant.
// Programss are expected to learn which indexes are useful for
// each operation.
// The int that each element refers to means something
// particular for each problem. Programs are also expected to
// learn what, of these, is important.
type I []*int

func (env *I) Copy() *I {
	if env == nil {
		return nil
	}
	newEnv := make(I, len(*env))
	for i, f := range *env {
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
func (env *I) New(envDiff []float64) *I {
	newEnv := make(I, len(envDiff))
	for i := range newEnv {
		newEnv[i] = new(int)
		var e int
		var f int
		if i < len(*env) {
			e = *(*env)[i]
		}
		if i < len(envDiff) {
			f = int(math.Ceil(envDiff[i]))
		}
		*newEnv[i] = e + f
	}
	//fmt.Println("Env.New length:", len(newEnv))
	return &newEnv
}

// Returns the absolute difference between the given environment
// and the passed in expectations. 0 in envDiff is treated as
// insinificant.
func (env *I) Diff(envDiff []float64) (diff int) {
	for i, f := range envDiff {
		if f != 0.0 {
			diff += int(math.Ceil(float64(*(*env)[i]) - f))
		}
	}
	return
}

func (env *I) MatchDiff(envDiff []float64) (diff int) {
	//fmt.Println("Env diff length:", len(*env), len(envDiff))
	for i, f := range envDiff {
		if i >= len(*env) {
			diff++
		} else if int(f) != *(*env)[i] {
			diff++
		}
	}
	return
}

func NewI(size int, baseVal int) *I {
	env := make(I, size)
	for i := 0; i < size; i++ {
		env[i] = new(int)
		*env[i] = baseVal
	}
	return &env
}
