package neural

import(
	"math"
	"math/rand"
)

func KeySet(m map[int]bool) []int {
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

// http://stackoverflow.com/questions/18390266/how-can-we-truncate-float64-type-to-a-particular-precision-in-golang
func round(num float64) int {
    return int(num + math.Copysign(0.5, num))
}
func toFixed(num float64, precision int) float64 {
    output := math.Pow(10, float64(precision))
    return float64(round(num * output)) / output
}