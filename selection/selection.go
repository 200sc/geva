package selection

import (
	"github.com/200sc/geva/pop"
	"math/rand"
)

// Return k random integers in range 0 to max
// Unknown author. Found here: http://play.golang.org/p/QH3_U3oiNL
func Sample(k, max int) (sampled []int) {
	swapped := make(map[int]int, k)
	for i := 0; i < k; i++ {
		// generate a random number r, where i <= r < max-min
		r := rand.Intn(max-i) + i

		// swapped[i], swapped[r] = swapped[r], swapped[i]
		vr, ok := swapped[r]
		if ok {
			sampled = append(sampled, vr)
		} else {
			sampled = append(sampled, r)
		}
		vi, ok := swapped[i]
		if ok {
			swapped[r] = vi
		} else {
			swapped[r] = i
		}
	}
	return
}

func KeySet_IntInt(m map[int]int) []int {
	keys := make([]int, len(m))

	i := 0
	for k := range m {
		keys[i] = k
		i++
	}

	return keys
}

type Method interface {
	Select(p *pop.Population) []pop.Individual
	GetParentProportion() int
}
