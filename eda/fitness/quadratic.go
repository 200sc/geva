package fitness

import "bitbucket.org/StephenPatrick/goevo/env"

// Quadratic is similar to Onemax, but specifically wants pairs of digits to be
// the same, and offers similar but not optimal fitness for pairs of 0s over
// pairs of 1s
func Quadratic(e *env.F) int {
	diff := float64(len(*e) / 2)
	for i := 0; i < len(*e)/2; i++ {
		diff -= quadf2(e.Get(2*i), e.Get((2*i)+1))
	}
	return int(diff)
}

// As defined in the BMDA paper
func quadf2(u, v float64) float64 {
	return .9 - (.9 * (u + v)) + (1.9 * u * v)
}
