package fitness

import (
	"math"

	"bitbucket.org/StephenPatrick/goevo/env"
)

// AlternatingABS is a fitness function which returns the
// best fitness for alternating ones and zeroes
func AlternatingABS(e *env.F) int {
	diff := 0.0
	for i := 0; i < len(*e); i++ {
		if i%2 == 0 {
			diff += math.Abs(e.Get(i) - 1)
		} else {
			diff += math.Abs(e.Get(i))
		}
	}
	return int(diff)
}
