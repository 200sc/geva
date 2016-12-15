package goevo

import (
	"math"
	"sort"
)

type TestCase struct {
	inputs  [][]float64
	outputs [][]float64
	title   string
}

func Pow8TestCase() TestCase {
	in := [][]float64{
		{1.0},
		{2.0},
		{3.0},
		{4.0},
	}

	out := make([][]float64, len(in))
	for i, f := range in {
		out[i] = []float64{math.Pow(f[0], 8.0)}
	}

	return TestCase{
		in,
		out,
		"Pow8",
	}
}

func PowSumTestCase() TestCase {
	in := [][]float64{
		{10, 1},
		//{10, 2},
		{20, 1},
		//{20, 2},
		{30, 1},
		//{30, 2},
	}

	out := make([][]float64, len(in))
	for i, f := range in {
		out[i] = []float64{PowSum(f[0], f[1])}
	}

	return TestCase{
		in,
		out,
		"PowSum",
	}
}

func PowSum(max, pow float64) float64 {
	out := 0.0
	for i := 0.0; i <= max; i++ {
		out += math.Pow(i, pow)
	}
	return out
}

func ReverseListTestCase() TestCase {
	in := [][]float64{
		{1.0, 2.0, 3.0, 4.0, 5.0},
		{7.0, 8.0, 9.0, 10.0, 11.0, 12.0},
		{15.0, 14.0, 13.0},
	}

	out := make([][]float64, len(in))
	for i, f := range in {
		out[i] = ReverseList(f)
	}

	return TestCase{
		in,
		out,
		"ReverseList",
	}
}

func ReverseList(lst []float64) []float64 {
	outList := make([]float64, len(lst))
	halflen := (len(lst) / 2) + 1
	for i := 0; i < halflen; i++ {
		outList[i] = lst[len(lst)-(i+1)]
	}
	return outList
}

func SortListTestCase() TestCase {
	in := [][]float64{
		{2.0, 3.0, 1.0, 5.0, 4.0},
		{7.0, 9.0, 8.0, 11.0, 12.0, 10.0},
		{15.0, 14.0, 13.0},
	}

	out := make([][]float64, len(in))
	for i, f := range in {
		out[i] = make([]float64, len(f))
		copy(out[i], f)
		sort.Float64s(out[i])
	}

	return TestCase{
		in,
		out,
		"SortList",
	}
}

func TransposeMatrixTestCase() TestCase {
	in := [][]float64{
		{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0},
	}

	out := [][]float64{
		{1.0, 4.0, 7.0, 2.0, 5.0, 8.0, 3.0, 6.0, 9.0},
	}
	return TestCase{
		in,
		out,
		"TransposeMatrix",
	}
}

func MultiplyMatrixTestCase() TestCase {
	in := [][]float64{
		{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0,
			3.0, 2.0, 0.0, 0.0, 1.0, 0.0, 4.0, 4.0, 3.0},
	}

	out := make([][]float64, len(in))
	for i, f := range in {
		l := len(f) / 2
		out[i] = MatrixMultiply(f[0:l], f[l:])
	}
	return TestCase{
		in,
		out,
		"MultiplyMatrix",
	}
}

// matrices are assumed to be square
func MatrixMultiply(a []float64, b []float64) []float64 {
	sqrtA := math.Sqrt(float64(len(a)))
	sqrtB := math.Sqrt(float64(len(b)))
	if sqrtA != math.Ceil(sqrtA) {
		panic("Bad input to Matrix Multiply: not a square matrix a")
	}
	if sqrtB != math.Ceil(sqrtB) {
		panic("Bad input to Matrix Multiply: not a square matrix b")
	}
	if sqrtA != sqrtB {
		panic("Non-multipliable matrices in Matrix Multiply")
	}
	out := make([]float64, len(a))
	// index
	for i := 0; i < len(a); i++ {
		// row to multiply
		k1 := i - i%int(sqrtA)
		// column to multiply
		j := i % int(sqrtA)
		for k := k1; k < k1+int(sqrtA); k++ {

			//fmt.Println("Adding to index", i, "k, j", k, j)
			out[i] += a[k] * b[j]

			j += int(sqrtA)
		}
	}
	return out
}
