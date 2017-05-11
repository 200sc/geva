package eda

import (
	"math"
	"sort"

	"bitbucket.org/StephenPatrick/goevo/env"
)

type MIMIC struct {
	Base
	Percentile float64
}

func (mimic *MIMIC) Adjust() Model {
	samples := NSamples(mimic.samples, mimic.F)
	fitnesses := SampleFitnesses(mimic, samples)

	// Filter the samples so that they are only those with a fitness under some
	// percentile of fitness
	thetaFitness := fitnesses[int(float64(mimic.samples)*mimic.Percentile)]
	filtered := []*env.F{}
	for _, s := range samples {
		mimic.F = s
		fitness := mimic.fitness(mimic)
		if fitness <= thetaFitness {
			filtered = append(filtered, s)
		}
	}
	// Recalculate mimic.F based on the filtered samples
	mimic.UpdateFromSamples(filtered)
	return mimic
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
	return fitnesses
}

func NSamples(n int, senv *env.F) []*env.F {
	samples := make([]*env.F, n)
	for i := 0; i < n; i++ {
		samples[i] = GetSample(senv)
	}
	return samples
}

func MIMICModel(opts ...Option) (Model, error) {
	var err error
	mimic := new(MIMIC)
	mimic.Percentile = .25
	mimic.Base, err = DefaultBase(opts...)
	// Generate a random population of samples
	samples := NSamples(mimic.samples, env.NewF(mimic.length, mimic.baseValue))
	// Get the median fitness of the sample set
	// fitnesses := SampleFitnesses(mimic, samples)
	// This seems useless so it is commented out
	mimic.UpdateFromSamples(samples)
	return mimic, err
}

func (mimic *MIMIC) UpdateFromSamples(samples []*env.F) {
	// Let mimic.F be the density estimator of the median fitness
	//
	// Find the element in the population with the lowest entropy
	minEntropyIndex, minF := MinEntropy(samples)
	*(*mimic.F)[minEntropyIndex] = minF

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
	prevIndex := minEntropyIndex
	for i := 1; i < mimic.length; i++ {
		index, f := MinConditionalEntropy(samples, prevIndex, &available)
		*(*mimic.F)[index] = f
		prevIndex = index
	}
}

func MinConditionalEntropy(samples []*env.F, prev int, available *[]int) (int, float64) {
	min := 0
	minV := math.MaxFloat64
	minF := 0.0
	minIndex := 0
	for j, i := range *available {
		hf, f := ConditionalEntropy(samples, i, prev)
		if hf < minV {
			minV = hf
			min = i
			minF = f
			minIndex = j
		}
	}
	// Remove the chosen index from the available list
	*available = append((*available)[:minIndex], (*available)[minIndex+1:]...)
	return min, minF
}

func MinEntropy(samples []*env.F) (int, float64) {
	min := 0
	minV := math.MaxFloat64
	minF := 0.0
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
			minF = f
		}
	}
	return min, minF
}

func Entropy(f float64) float64 {
	if f == 0 || f == 1 {
		return 0
	}
	return -1 * ((f * math.Log2(f)) + ((1 - f) * math.Log2(1-f)))
}

func ConditionalEntropy(samples []*env.F, a, b int) (float64, float64) {
	// Assume bitstrings, assumption can be removed later
	h := 0.0
	patb := 0.0
	// Iterate over all potential values for env[a] and env[b]
	for acheck := 0.0; acheck <= 1.0; acheck++ {
		for bcheck := 0.0; bcheck <= 1.0; bcheck++ {
			// Find the total number of times in the samples where a and b are
			// the expected value pair
			// and where b is the expected value
			totab := 0
			totb := 0
			for _, s := range samples {
				af := *(*s)[a]
				bf := *(*s)[b]
				if bf == bcheck {
					if af == acheck {
						totab++
					}
					totb++
				}
			}
			// Divide by total samples to get probabilities
			pab := float64(totab) / float64(len(samples))
			pb := float64(totb) / float64(len(samples))
			if acheck == 1.0 {
				patb += pab
			}
			// Avoid impossible math
			if pab == 0.0 || pb == 0.0 {
				continue
			}
			h += pab * math.Log2(pb/pab)
		}
	}
	// Also need to return P(a=T|b)
	return h, patb
}
