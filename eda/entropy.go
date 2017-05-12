package eda

import (
	"math"

	"bitbucket.org/StephenPatrick/goevo/env"
)

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
	// Also need to return P(a=T|b)
	return h
}
