package neural

func KeySet(m map[int]bool) []int {
	keys := make([]int, len(m))

	i := 0
	for k := range m {
	    keys[i] = k
	    i++
	}

	return keys
}