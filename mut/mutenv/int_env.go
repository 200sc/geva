package mutenv

import (
	"github.com/200sc/geva/env"
	"github.com/200sc/geva/mut"
)

type I func(*env.I)

func OnAll(f mut.FloatMutator) I {
	return func(e *env.I) {
		for i := 0; i < len(*e); i++ {
			*(*e)[i] = int(f(float64(*(*e)[i])))
		}
	}
}

func And(f1, f2 I) I {
	return func(e *env.I) {
		f1(e)
		f2(e)
	}
}
