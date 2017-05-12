package eda

import (
	"sort"

	"bitbucket.org/StephenPatrick/goevo/env"
)

type MIMIC struct {
	Base
	PTF     *env.F
	Indices []int
}

func (mimic *MIMIC) Adjust() Model {
	samples := mimic.NSamples()
	fitnesses := SampleFitnesses(mimic, samples)

	// Filter the samples so that they are only those with a fitness under some
	// percentile of fitness
	thetaFitness := fitnesses[int(float64(mimic.samples)*mimic.learningRate)]
	filtered := []*env.F{}
	for _, s := range samples {
		mimic.F = s
		fitness := mimic.fitness(mimic)
		if fitness <= thetaFitness {
			filtered = append(filtered, s)
		}
	}
	//fmt.Println("Length of filtered", len(filtered))
	// Recalculate mimic.F based on the filtered samples
	mimic.UpdateFromSamples(filtered)
	mimic.F.Mutate(mimic.mutationRate, mimic.fmutator)
	mimic.PTF.Mutate(mimic.mutationRate, mimic.fmutator)
	mimic.learningRate = mimic.lmutator(mimic.learningRate)
	return mimic
}

func (mimic *MIMIC) GetSample() *env.F {
	// A mimic sample goes through mimic.Indices
	s := env.NewF(mimic.length, 0.0)
	// Index zero is univariate, stored in the PTT environment
	s.Set(mimic.Indices[0], mimic.F.GetBinary(mimic.Indices[0]))
	// Each following index is based on whatever exists in the previous index
	for i := 1; i < len(mimic.Indices); i++ {
		var e *env.F
		if s.Get(mimic.Indices[i-1]) == 0.0 {
			e = mimic.PTF
		} else {
			e = mimic.F
		}
		s.Set(mimic.Indices[i], e.GetBinary(mimic.Indices[i]))
	}
	return s
}

func (mimic *MIMIC) NSamples() []*env.F {
	samples := make([]*env.F, mimic.samples)
	for i := 0; i < mimic.samples; i++ {
		samples[i] = mimic.GetSample()
	}
	return samples
}

func SampleFitnesses(m Model, samples []*env.F) []int {
	bm := m.BaseModel()
	initF := bm.F.Copy()
	fitnesses := make([]int, len(samples))
	for i, s := range samples {
		bm.F = s
		fitnesses[i] = bm.fitness(m)
	}
	sort.Slice(fitnesses, func(i, j int) bool { return fitnesses[i] < fitnesses[j] })
	bm.F = initF
	//fmt.Println(fitnesses)
	return fitnesses
}

func MIMICModel(opts ...Option) (Model, error) {
	var err error
	mimic := new(MIMIC)
	mimic.Base, err = DefaultBase(opts...)
	mimic.PTF = mimic.F.Copy()
	// Generate a random population of samples
	samples := NSamples(mimic.samples, env.NewF(mimic.length, mimic.baseValue))
	// Get the median fitness of the sample set
	// fitnesses := SampleFitnesses(mimic, samples)
	// This seems useless so it is commented out
	mimic.Indices = make([]int, mimic.length)
	mimic.UpdateFromSamples(samples)
	return mimic, err
}

func (mimic *MIMIC) UpdateFromSamples(samples []*env.F) {
	// Let mimic.F be the density estimator of the median fitness
	//
	// Find the element in the population with the lowest entropy
	minEntropyIndex, minF := MinEntropy(samples)
	*(*mimic.F)[minEntropyIndex] = minF
	mimic.Indices[0] = minEntropyIndex

	// Remaining indicies
	available := make([]int, mimic.length)
	for i := range available {
		available[i] = i
	}
	// Remove the initial index from the available list of indices
	available = append(available[:minEntropyIndex], available[minEntropyIndex+1:]...)
	// For each following element, find the element in the population
	// where the entropy of the element is minimized, given the previous
	// element.
	for i := 1; i < mimic.length; i++ {
		index := MinConditionalEntropy(samples, mimic.Indices[i-1], &available)
		ptt, ptf := BitStringBivariate(samples, index, mimic.Indices[i-1])
		*(*mimic.F)[index] = ptt
		*(*mimic.PTF)[index] = ptf
		mimic.Indices[i] = index
	}
	//fmt.Println(mimic.PTF)
}