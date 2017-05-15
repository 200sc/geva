package stat

var (
	factMemo = map[int]int{
		1:  1,
		0:  1,
		-1: 0,
	}
)

func Factorial(n int) int {
	if _, ok := factMemo[n]; !ok {
		f := Factorial(n - 1)
		factMemo[n] = n * f
	}
	return factMemo[n]
}
