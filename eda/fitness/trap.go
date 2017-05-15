package fitness

import (
	"math/rand"

	"bitbucket.org/StephenPatrick/goevo/env"
)

// TrapABS takes in the bitlength t of the trap
// trap fitnesses follow reward worse fitness for
// increasing numbers of ones in subgroups of the environment
// until all elements in the subgroup are one, at which
// point the reward the best fitness.
// A graph:
//
// Good fitness  |                             -
//               |--------
//               |        -----------
//               |                   ----------
//               |------------------------------
// Poor Fitness   No ones              More ones
//
func TrapABS(t int) func(e *env.F) int {
	return func(e *env.F) int {
		diff := 0
		for i := 0; (i + t) < len(*e); i += t {
			ones := 0
			for j := i; j < i+t; j++ {
				if rand.Float64() < e.Get(j) {
					ones++
				}
			}
			diff += trap(t, ones)
		}
		return diff
	}
}

func trap(t, ones int) int {
	if t == ones {
		return 0
	}
	return ones + 1
}
