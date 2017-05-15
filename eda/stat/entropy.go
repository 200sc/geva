package stat

import (
	"math"

	"bitbucket.org/StephenPatrick/goevo/env"
)

// MinConditionalEntropy finds the index of the given samples with the
// lowest entropy given the provided previous index, only checking indices
// in the available list. Afterward, the index chosen is removed from the
// available list.
func MinConditionalEntropy(samples []*env.F, prev int, available *[]int) int {
	min := 0
	minV := math.MaxFloat64
	minIndex := 0
	for j, i := range *available {
		hf := ConditionalEntropy(samples, i, prev)
		if hf < minV {
			minV = hf
			min = i
			minIndex = j
		}
	}
	// Remove the chosen index from the available list
	*available = append((*available)[:minIndex], (*available)[minIndex+1:]...)
	return min
}

// MinEntropy finds the index of the given samples with the lowest
// entropy. Also returns the univariate probability of the returned index.
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

// Entropy returns the entropy of a float64
func Entropy(f float64) float64 {
	if f == 0 || f == 1 {
		return 0
	}
	return -1 * ((f * math.Log2(f)) + ((1 - f) * math.Log2(1-f)))
}

// ConditionalEntropy returns the conditional entropy of (a|b)
// in the assumption that the provided samples are bitstrings
func ConditionalEntropy(samples []*env.F, a, b int) float64 {
	// Assume bitstrings, assumption can be removed later
	h := 0.0
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
			// Avoid impossible math
			if pab == 0.0 || pb == 0.0 {
				continue
			}
			h += pab * math.Log2(pb/pab)
		}
	}
	return h
}

// MarginalEntropy returns the marginal entropy of a subset of the given
// samples, where the subset is defined by a slice of indices in the sample
// environments.
func MarginalEntropy(samples []*env.F, indices []int) float64 {
	combs := bsComb(len(indices))
	h := 0.0
	for _, c := range combs {
		occurrences := combOccurrences(samples, c, indices)
		if occurrences == 0 {
			continue
		}
		chance := float64(occurrences) / float64(len(samples))
		// It's not clear to me what we should be dividing chance by.
		// The first paper does not go into detail, and just says to obtain
		// the marginal entropy.
		//
		// The second paper says that chance should be dividing N_p, but
		// defines N_p as the total number of individuals in the population,
		// which is going to be HUGE, and shouldn't factor into the entropy
		// of a subset of a population?
		//
		// Using 1, this is joint entropy.
		h += chance * math.Log2(1/chance)
	}
	return h
}

// combOccurrences finds the number of times combination comb occurs in the
// sample set, where the indices in the samples of comb are defined by indices
func combOccurrences(samples []*env.F, comb []float64, indices []int) (occurrences int) {
sampleLoop:
	for _, s := range samples {
		for i, j := range indices {
			if s.Get(j) != comb[i] {
				continue sampleLoop
			}
		}
		occurrences++
	}
	return
}

// bsComb returns all combinations of bitstrings that can appear
// given indices number of bits
func bsComb(indices int) [][]float64 {
	if indices == 0 {
		return make([][]float64, 0)
	}
	combs := make([][]float64, 2)
	combs[0] = []float64{1.0}
	combs[1] = []float64{0.0}
	for i := 1; i < indices; i++ {
		newCombs := make([][]float64, 2*len(combs))
		for j, c := range combs {
			newCombs[2*j] = append(c, 1.0)
			newCombs[(2*j)+1] = append(c, 0.0)
		}
		combs = newCombs
	}
	return combs
}
