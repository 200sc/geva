package env

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
// insinificant.
func (env *F) Diff(envDiff []float64) (diff int) {
	diff_f := 0.0
	for i, f := range envDiff {
		if f != 0.0 {
			diff_f += *(*env)[i] - f
		}
	}
	diff = int(diff_f)
	return
}

func (env *F) SetAll(f float64) {
	for i := range *env {
		*(*env)[i] = f
	}
}

func NewF(size int, baseVal int) *F {
	env := make(F, size)
	baseVal_f := float64(baseVal)
	for i := 0; i < size; i++ {
		env[i] = new(float64)
		*env[i] = baseVal_f
	}
	return &env
}
