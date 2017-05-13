package env

import (
	//"fmt"
	"math"
	"strconv"
)

// I represents an integer-valued environment
type I []*int

// Copy copies an integer environment
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

// New Creates a combined environment by the given envDiff inputs
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

// Diff returns the absolute difference between the given environment
// and the passed in expectations. 0 in envDiff is treated as
// insignificant.
func (env *I) Diff(envDiff []float64) (diff int) {
	for i, f := range envDiff {
		if f != 0.0 {
			diff += int(math.Ceil(math.Abs(float64(*(*env)[i]) - f)))
		}
	}
	return
}

// MatchDiff compares each element in envDiff and env and
// retuns the number of elements which are not the same.
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

// NewI constructs an I with a baseValue at each space
func NewI(size int, baseVal int) *I {
	env := make(I, size)
	for i := 0; i < size; i++ {
		env[i] = new(int)
		*env[i] = baseVal
	}
	return &env
}

func (env *I) String() string {
	str := "["
	for _, v := range *env {
		str += strconv.Itoa(*v)
		str += " "
	}
	str += "]"
	return str
}
