package goevo

import (
	"math"
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
