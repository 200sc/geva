package eda

import (
	"math/rand"
	"sort"

	"bitbucket.org/StephenPatrick/goevo/env"
	"bitbucket.org/StephenPatrick/goevo/pop"
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

// PopEnvs converts a population of EnvInd members to a sample array
func PopEnvs(p *pop.Population) []*env.F {
	return MemberEnvs(p.Members)
}

// MemberEnvs converts a set of EnvInd members to a smaple array
func MemberEnvs(mem []pop.Individual) []*env.F {
	envs := make([]*env.F, len(mem))
	for i, s := range mem {
		envs[i] = s.(*EnvInd).F
	}
	return envs
}

// ReplaceLowFitnesses sorts the given population members by fitness
// and replaces the lowest members of those with the newMembers.
func (b *Base) ReplaceLowFitnesses(p *pop.Population, newMembers []*env.F) {
	envs := PopEnvs(p)

	// Sort envs by fitness
	_, envs = SampleFitnesses(b, envs)
	i := 0
	for j := len(envs) - 1; j >= 0; j-- {
		envs[j] = newMembers[i]
		i++
		if i >= len(newMembers) {
			break
		}
	}

	// set the population to be envs cast back to members
	for i := range p.Members {
		p.Members[i] = &EnvInd{envs[i]}
	}
}
