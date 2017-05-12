package eda

import "bitbucket.org/StephenPatrick/goevo/env"

// MIMIC is an EDA of the Mutual information maximizing input clustering algorithm
type MIMIC struct {
	Base
	PTF     *env.F
	Indices []int
}

// Adjust on a MIMIC generates samples from MIMIC's chain model, filters
// them with a straight greed selection algorithm to get the top-percentile
// samples and retrains its chain on the top-percentile samples.
// Todo: would mimic perform way better if it wasn't using a straight greed
// selection? (the paper doesn't refer to this selection algorithm as a selection
// algorithm but it totally is)
func (mimic *MIMIC) Adjust() Model {
	fitnesses, samples := SampleFitnesses(mimic, mimic.NSamples())

	// Filter the samples so that they are only those with a fitness under some
	// percentile of fitness
	thetaFitness := fitnesses[int(float64(mimic.samples)*mimic.learningRate)]
	fi := len(samples)
	for i, f := range fitnesses {
		if f > thetaFitness {
			fi = i
			break
		}
	}
	filtered := samples[0:fi]
	// Recalculate mimic.F based on the filtered samples
	// In other models we'd reset mimic.F here but we don't need to do that,
	// as we trash our samples and our old F both
	mimic.UpdateFromSamples(filtered)
	mimic.PTF.Mutate(mimic.mutationRate, mimic.fmutator)
	return mimic
}

// GetSample on a MIMIC iterates through the marked indices of the model
// where the first index uses a univariate sampling and each following index
// is a bivariate sampling based on the result of the previous sampled index.
func (mimic *MIMIC) GetSample() *env.F {
	// A mimic sample goes through mimic.Indices
	s := env.NewF(mimic.length, 0.0)

	// Index zero is univariate, stored in the PTT environment
	s.Set(mimic.Indices[0], mimic.F.GetBinary(mimic.Indices[0]))

	// Each following index is based on whatever exists in the previous index
	for i := 1; i < len(mimic.Indices); i++ {
		e := ConditionedBSEnv(s, mimic.F, mimic.PTF, mimic.Indices[i-1])
		s.Set(mimic.Indices[i], e.GetBinary(mimic.Indices[i]))
	}
	return s
}

// NSamples runs mimic.GetSample n times to produce a sample list
func (mimic *MIMIC) NSamples() []*env.F {
	samples := make([]*env.F, mimic.samples)
	for i := 0; i < mimic.samples; i++ {
		samples[i] = mimic.GetSample()
	}
	return samples
}

// MIMICModel returns an initialized MIMIC EDA
func MIMICModel(opts ...Option) (Model, error) {
	var err error
	mimic := new(MIMIC)
	mimic.Base, err = DefaultBase(opts...)
	// We initialize with -1 so that if something doesn't get replaced
	// due to an issue with the algorithm we will crash sooner, or can
	// potentially notice it.
	mimic.PTF = env.NewF(mimic.length, -1.0)
	// Generate a random population of samples
	samples := NSamples(mimic.samples, env.NewF(mimic.length, mimic.baseValue))
	// Get the median fitness of the sample set
	// fitnesses := SampleFitnesses(mimic, samples)
	// This seems useless so it is commented out
	mimic.Indices = make([]int, mimic.length)
	mimic.UpdateFromSamples(samples)
	return mimic, err
}

// UpdateFromSamples updates the two floating point vectors that
// a MIMIC stores, the former for the probability that an element is
// true given the former in the index list is true, and the latter
// given the former in the index list is false.
func (mimic *MIMIC) UpdateFromSamples(samples []*env.F) {
	// Let mimic.F be the density estimator of the median fitness
	//
	// Find the element in the population with the lowest entropy
	minEntropyIndex, minF := MinEntropy(samples)
	*(*mimic.F)[minEntropyIndex] = minF
	*(*mimic.PTF)[minEntropyIndex] = minF
	mimic.Indices[0] = minEntropyIndex

	// Remaining indicies
	available := mimic.GenIndices()
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
	//time.Sleep(1 * time.Second)
}
