package eda

import (
	"math"
	"sort"

	"bitbucket.org/StephenPatrick/goevo/env"
)

type MIMIC struct {
	Base
	MedianFitness int
}

func (mimic *MIMIC) Adjust() Model {

	return mimic
}

func MIMICModel(opts ...Option) (Model, error) {
	var err error
	mimic := new(MIMIC)
	mimic.Base, err = DefaultBase(opts...)
	senv := env.NewF(mimic.length, mimic.baseValue)
	// Generate a random population of samples
	samples := make([]*env.F, mimic.samples)
	for i := 0; i < mimic.samples; i++ {
		samples[i] = GetSample(senv)
	}
	// Get the median fitness of the sample set
	fitnesses := make([]int, mimic.samples)
	for i, s := range samples {
		mimic.F = s
		fitnesses[i] = mimic.fitness(mimic)
	}
	sort.Slice(fitnesses, func(i, j int) bool { return fitnesses[i] < fitnesses[j] })
	mimic.MedianFitness = fitnesses[mimic.samples/2]
	// Let mimic.F be the density estimator of the median fitness
	//
	// Find the element in the population with the lowest entropy
	// Set that element to be mimic.F[0]
	indices := make([]int, mimic.length)
	minEntropyIndex := MinEntropy(samples)
	indices[0] = minEntropyIndex

	// For each following element, find the element in the population
	// where the entropy of the element is minimized, given the previous
	// element.
	for i := 1; i < mimic.length; i++ {
		indices[i] = MinEntropyGiven(samples, indices[i-1])
	}

	// Somehow turn these indices into something we can get the fitness of
	return mimic, err
}

func MinEntropyGiven(samples []*env.F, prev int) int {

}

func MinEntropy(samples []*env.F) int {
	min := 0
	minV := math.MaxFloat64
	for i := 0; i < len(*samples[0]); i++ {
		f := 0.0
		for _, s := range samples {
			f += *(*s)[i]
		}
		f /= float64(len(samples))
		hf := Entropy(f)
		if hf < minV {
			minV = hf
			min = i
		}
	}
	return min
}

func Entropy(f float64) float64 {
	if f == 0 || f == 1 {
		return 0
	}
	return -1 * ((f * math.Log2(f)) + ((1 - f) * math.Log2(1-f)))
}
