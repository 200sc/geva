package alg

import (
	"math/rand"
)

// AKA Single-Instance Roulette Search
// Given a set of cumulative weights, pick an index
// evenly distributed as according to said weights.
// Tl;dr Binary search but with ranges of values
func CumWeightedChooseOne(remainingWeights []float64) int {
	totalWeight := remainingWeights[len(remainingWeights)-1]
	choice := rand.Float64() * totalWeight
	i := len(remainingWeights) / 2
	start := 0
	end := len(remainingWeights) - 1
	for {
		if remainingWeights[i] >= choice {
			if i != 0 && remainingWeights[i-1] > choice {
				end = i
				i = (start + end) / 2
			} else {
				return i
			}
		} else if remainingWeights[i] < choice {
			if remainingWeights[i+1] < choice {
				start = i
				i = (start + end) / 2
			} else {
				return i + 1
			}
		}
	}
}
