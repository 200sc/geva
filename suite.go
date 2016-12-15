package goevo

import (
	"fmt"
	"goevo/alg"
	"goevo/pop"
	"math"
	//"runtime"
	"time"
)

type SuiteFunc func(interface{}, int) []pop.Individual

func RunSuite(testCases []TestCase, demeCount, popSize, testGenerations int, options interface{},
	suiteFunc SuiteFunc, selection []pop.SMethod, pairing []pop.PMethod,
	goal int, elites alg.IntRange, migration float64) {

	for _, tc := range testCases {
		totalGenerations := 0
		fmt.Println(tc.title)
		loops := 1
		mean := 0.0
		nextPrint := 5000
		results := []float64{}
		timings := []time.Duration{}
		fitnesses := []int{}
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

			// Todo: Think about what kind of memory
			// statistics make sense to look into
			// var mem runtime.MemStats
			// runtime.ReadMemStats(&mem)
			// fmt.Println(mem.Alloc)
			// fmt.Println(mem.TotalAlloc)
			// fmt.Println(mem.HeapAlloc)
			// fmt.Println(mem.HeapSys)
			t1 := time.Now()
			for j := 0; j < numGens; j++ {
				stopEarly := dg.NextGeneration()
				if j == numGens-1 || stopEarly {
					t2 := time.Since(t1) / time.Duration(j+1)
					timings = append(timings, t2)

					totalGenerations += j + 1

					results = append(results, float64(j+1))

					if loops%20 == 1 || totalGenerations > nextPrint {
						fmt.Println("Loop", loops, "Gens", totalGenerations)
						ind, fitness := dg.BestMember()
						ind.Print()
						mean = float64(totalGenerations) / float64(loops)
						fmt.Println("Best fitness reached: ", fitness)
						fitnesses = append(fitnesses, fitness)
						fmt.Println("Generations taken: ", j+1)
						fmt.Println("Average Generations: ", mean)
						fmt.Println("Time taken per generation:", t2)
						nextPrint = totalGenerations + 5000
					}
					break
				}
			}
			loops++
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

		timeMean := time.Duration(0)
		for _, f := range timings {
			timeMean += f
		}
		timeMean /= time.Duration(len(timings))

		timeStdvTotal := 0.0
		for _, f := range timings {
			timeStdvTotal += math.Pow(float64(f-timeMean), 2)
		}
		timeStdvTotal /= float64(len(timings))
		timeStdv := math.Sqrt(timeStdvTotal)
		fmt.Println("Average time per generation:", timeMean)
		fmt.Println("Time per generation Standard Deviation:", time.Duration(int(timeStdv)))

		fitnessMean := 0
		for _, f := range fitnesses {
			fitnessMean += f
		}
		fmt.Println(fitnessMean, len(fitnesses))
		fitnessMean /= len(fitnesses)

		fitnessStdv := 0.0
		for _, f := range fitnesses {
			fitnessStdv += math.Pow(float64(f-fitnessMean), 2)
		}
		fitnessStdv /= float64(len(fitnesses))
		fitnessStdv = math.Sqrt(fitnessStdv)

		fmt.Println("Average best fitness:", fitnessMean)
		fmt.Println("Stdv best fitness:", fitnessStdv)
	}
}
