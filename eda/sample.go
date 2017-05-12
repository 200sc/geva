package eda

import (
	"math/rand"
	"sort"

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

// NSamples returns an array of n samples from calling GetSample
func NSamples(n int, senv *env.F) []*env.F {
	samples := make([]*env.F, n)
	for i := 0; i < n; i++ {
		samples[i] = GetSample(senv)
	}
	return samples
}

// UnivariateFromSamples returns the average value within a given
// index in a sample set. Used in BMDA for roots.
func UnivariateFromSamples(samples []*env.F, a int) float64 {
	f := 0.0
	for _, s := range samples {
		f += s.Get(a)
	}
	f /= float64(len(samples))
	return f
}

// SampleFitnesses returns a sorted list of fitnesses for the given
// samples in model m, and sorts samples so that the fitness of samples[i] is
// in fitnesses[i]
func SampleFitnesses(m Model, samples []*env.F) []int {
	bm := m.BaseModel()
	initF := bm.F.Copy()
	fitnesses := make([]int, len(samples))
	for i, s := range samples {
		bm.F = s
		fitnesses[i] = bm.fitness(m)
	}
	sort.Slice(fitnesses, func(i, j int) bool { return fitnesses[i] < fitnesses[j] })
	sort.Slice(samples, func(i, j int) bool { return fitnesses[i] < fitnesses[j] })
	bm.F = initF
	//fmt.Println(fitnesses)
	return fitnesses
}
