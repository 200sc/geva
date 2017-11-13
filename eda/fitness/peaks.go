package fitness

import (
	"math"
	"math/rand"

	"github.com/200sc/geva/env"
)

// FourPeaks represents a problem where there are four explicit
// maxima in the search space and two of the maxima can hide the
// other two.
func FourPeaks(t int) func(e *env.F) int {
	return func(e *env.F) int {
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

func SixPeaks(t int) func(e *env.F) int {
	return func(e *env.F) int {
		leadingOnes, leadingZeroes, trailingOnes, trailingZeroes := bsEndlengths(e)
		base := int(math.Max(float64(leadingOnes), float64(trailingZeroes)))
		if (trailingZeroes > t && leadingOnes > t) ||
			(trailingOnes > t && leadingZeroes > t) {
			base += len(*e)
		}
		return ((2 * len(*e)) - t) - base
	}
}

func bsEndlengths(e *env.F) (int, int, int, int) {
	leadingOnes := 0
	leadingZeroes := 0
	if rand.Float64() < e.Get(0) {
		leadingOnes++
		for i := 1; i < len(*e); i++ {
			if rand.Float64() < e.Get(i) {
				leadingOnes++
			} else {
				break
			}
		}
	} else {
		leadingZeroes++
		for i := 1; i < len(*e); i++ {
			if rand.Float64() > e.Get(i) {
				leadingZeroes++
			} else {
				break
			}
		}
	}
	trailingOnes := 0
	trailingZeroes := 0
	if rand.Float64() < e.Get(len(*e)-1) {
		trailingOnes++
		for i := len(*e) - 2; i > -1; i-- {
			if rand.Float64() < e.Get(i) {
				trailingOnes++
			} else {
				break
			}
		}
	} else {
		trailingZeroes++
		for i := len(*e) - 2; i > -1; i-- {
			if rand.Float64() > e.Get(i) {
				trailingZeroes++
			} else {
				break
			}
		}
	}
	return leadingOnes, leadingZeroes, trailingOnes, trailingZeroes
}

// Related problems:
// K-Coloring
