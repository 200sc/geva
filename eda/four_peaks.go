package eda

import (
	"math"
	"math/rand"
)

// FourPeaks represents a problem where there are four explicit
// maxima in the search space and two of the maxima can hide the
// other two.
func FourPeaks(t int) func(m Model) int {
	return func(m Model) int {
		e := m.ToEnv()
		leadingOnes := 0
		for _, f := range *e {
			if rand.Float64() < *f {
				leadingOnes++
			} else {
				break
			}
		}
		trailingZeroes := 0
		for i := len(*e) - 1; i >= 0; i-- {
			f := (*e)[i]
			if rand.Float64() > *f {
				trailingZeroes++
			} else {
				break
			}
		}
		base := int(math.Max(float64(leadingOnes), float64(trailingZeroes)))
		if trailingZeroes > t && leadingOnes > t {
			base += len(*e)
		}
		return ((2 * len(*e)) - t) - base
	}
}

// Related problems:
// Six peaks
// K-Coloring
