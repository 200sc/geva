package goevo

import (
	"fmt"
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

	out := [][]float64{
		{math.Pow(1.0, 8.0)},
		{math.Pow(2.0, 8.0)},
		{math.Pow(3.0, 8.0)},
		{math.Pow(4.0, 8.0)},
	}
	title := "Pow8"
	return GPTestCase{
		in,
		out,
		title,
	}
}

type SuiteFunc func(interface{}, int) []population.Individual

func RunSuite(testCases []GPTestCase, demeCount, popSize, testGenerations int, options interface{},
	suiteFunc SuiteFunc, selection population.SelectionMethod, pairing population.PairingMethod,
	elites, goal int, migration float64) {

	for _, tc := range testCases {
		totalGenerations := 0
		fmt.Println(tc.title)
		loops := 0
		variance := 0.0
		oldMean := 0.0
		mean := 0.0
		stdDev := 0.0
		for totalGenerations < testGenerations {

			dg := MakeDemes(
				demeCount,
				suiteFunc(options, popSize),
				selection,
				pairing,
				tc.inputs,
				tc.outputs,
				len(tc.inputs),
				elites,
				goal,
				migration)

			numGens := 5000

			for j := 0; j < numGens; j++ {
				stopEarly := dg.NextGeneration()
				if j == numGens-1 || stopEarly {
					totalGenerations += j + 1

					if loops%20 == 1 {
						fmt.Println("Loop", loops, "Gens", totalGenerations)
						ind, _ := dg.BestMember()
						ind.Print()
						oldMean = mean
						mean = float64(totalGenerations) / float64(loops)
						variance = variance + (float64(j+1)-oldMean)*(float64(j+1)-mean)
						stdDev = math.Sqrt(variance / float64(totalGenerations-1))
						fmt.Println("Generations taken: ", j+1)
						fmt.Println("Average Generations: ", mean)
						fmt.Println("Standard Deviation: ", stdDev)
					}
					break
				}
			}
			loops += 1
		}
		oldMean = mean
		mean = float64(totalGenerations) / float64(loops)
		fmt.Println("End Average Generations: ", mean)
		fmt.Println("End Standard Deviation: ", stdDev)
	}
}
