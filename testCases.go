package goevo

import (
	"fmt"
	"goevo/alg"
	"goevo/population"
	"math"
)

type GPTestCase struct {
	inputs  [][]float64
	outputs [][]float64
	title   string
}

func Pow8TestCase() GPTestCase {
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

	return GPTestCase{
		in,
		out,
		"Pow8",
	}
}

func PowSumTestCase() GPTestCase {
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

	return GPTestCase{
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

func ReverseListTestCase() GPTestCase {
	in := [][]float64{
		{1.0, 2.0, 3.0, 4.0, 5.0},
		{7.0, 8.0, 9.0, 10.0, 11.0, 12.0},
		{15.0, 14.0, 13.0},
	}

	out := make([][]float64, len(in))
	for i, f := range in {
		out[i] = ReverseList(f)
	}

	return GPTestCase{
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

type SuiteFunc func(interface{}, int) []population.Individual

func RunSuite(testCases []GPTestCase, demeCount, popSize, testGenerations int, options interface{},
	suiteFunc SuiteFunc, selection []population.SelectionMethod, pairing []population.PairingMethod,
	goal int, elites alg.IntRange, migration float64) {

	for _, tc := range testCases {
		totalGenerations := 0
		fmt.Println(tc.title)
		loops := 1
		mean := 0.0
		nextPrint := 5000
		results := []float64{}
		for totalGenerations < testGenerations {

			dg := MakeDemes(
				demeCount,
				suiteFunc(options, popSize),
				selection,
				pairing,
				tc.inputs,
				tc.outputs,
				len(tc.inputs),
				goal,
				elites,
				migration)

			numGens := 5000

			for j := 0; j < numGens; j++ {
				stopEarly := dg.NextGeneration()
				if j == numGens-1 || stopEarly {
					totalGenerations += j + 1

					if loops%20 == 1 || totalGenerations > nextPrint {
						fmt.Println("Loop", loops, "Gens", totalGenerations)
						ind, _ := dg.BestMember()
						ind.Print()
						mean = float64(totalGenerations) / float64(loops)
						results = append(results, float64(j+1))
						fmt.Println("Generations taken: ", j+1)
						fmt.Println("Average Generations: ", mean)
						nextPrint = totalGenerations + 5000
					}
					break
				}
			}
			loops += 1
		}
		mean = float64(totalGenerations) / float64(loops)
		stdDevTotal := 0.0
		for _, f := range results {
			stdDevTotal += math.Pow((f - mean), 2)
		}
		stdDevTotal /= float64(len(results))
		stdDev := math.Sqrt(stdDevTotal)
		fmt.Println("End Average Generations: ", mean)
		fmt.Println("End Standard Deviation: ", stdDev)
	}
}
