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
func SampleFitnesses(m Model, samples []*env.F) ([]int, []*env.F) {
	bm := m.BaseModel()
	fitSamples := make([]FitSample, len(samples))
	for i, s := range samples {
		fitSamples[i] = FitSample{s, bm.fitness(s)}
	}
	sort.Slice(fitSamples, func(i, j int) bool {
		return fitSamples[i].fitness < fitSamples[j].fitness
	})
	fitnesses := make([]int, len(samples))
	for i := range fitSamples {
		fitnesses[i] = fitSamples[i].fitness
		samples[i] = fitSamples[i].sample
	}
	return fitnesses, samples
}

// FitSample is a tuple of a sample and its fitness
type FitSample struct {
	sample  *env.F
	fitness int
}
