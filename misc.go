package neural

import(
	"math"
)

func KeySet(m map[int]interface{}) []int {
	keys := make([]int, len(m))

	i := 0
	for k := range m {
	    keys[i] = k
	    i++
	}

	return keys
}

func RefreshingRand(ch chan<- float64) {
	for {
		// We refresh the random number here--
		// after about 40 times we would expect
		// the number to start running out of digits
		next := rand.Float64()
		for i := 0; i < 40; i++ {
			ch <- next
			next = next * 10
			next = next - (math.Floor(next))
		}
	}
}