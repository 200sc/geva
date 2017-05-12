package eda

import (
	"math/rand"

	"bitbucket.org/StephenPatrick/goevo/env"
)

// GetSample and NSamples both are univariate

// GetSample returns an environment candidate where each input element is treated
// as a percent from 0 to 1 inclusive, and each output is each input randomized to
// either 1 or 0.
func GetSample(e *env.F) *env.F {
	sample := e.Copy()
	for _, f := range *sample {
		if rand.Float64() <= *f {
			*f = 1
		} else {
			*f = 0
		}
	}
	return sample
}

func NSamples(n int, senv *env.F) []*env.F {
	samples := make([]*env.F, n)
	for i := 0; i < n; i++ {
		samples[i] = GetSample(senv)
	}
	return samples
}
