package goevo

import (
	"encoding/csv"
	"fmt"
	"bitbucket.org/StephenPatrick/goevo/alg"
	"bitbucket.org/StephenPatrick/goevo/pop"
	"math"
	"os"
	"sort"
	"strconv"
	"time"
)

type TestSuite struct {
	testCases                           []TestCase
	demeCount, popSize, testGenerations int
	options                             interface{}
	suiteFunc                           SuiteFunc
	initOptions                         interface{}
	suiteInitFunc                       SuiteInitFunc
	selection                           []pop.SMethod
	pairing                             []pop.PMethod
	goal                                int
	elites                              alg.IntRange
	migration                           float64
	titleSuffix                         string
}

func RunTestSuites(suites []TestSuite) {
	writer := SuiteWriter()

	for _, suite := range suites {
		RunSingleSuite(suite, writer)
		writer.Flush()
	}

	writer.Flush()
}

type SuiteFunc func(interface{}, int) []pop.Individual

type SuiteInitFunc func(interface{})

func SuiteWriter() *csv.Writer {
	file := "logs/log"
	file += time.Now().Format("_Jan_2_15-04-05_2006")
	file += ".txt"
	fHandle, _ := os.Create(file)
	writer := csv.NewWriter(fHandle)

	writer.Write([]string{
		"Test Title",
		"Total Generations",
		"Mean Generations/Solution",
		"Stdv Generations/Solution",
		"Total Time",
		"Mean Time/Gen",
		"Stdv Time/Gen",
		"Total Best Fitness/Gen",
		"Mean Best Fitness/Gen",
		"Stdv Best Fitness/Gen",
		"Total Average Fitness/Gen",
		"Mean Average Fitness/Gen",
		"Stdv Average Fitness/Gen",
	})
	return writer
}

func RunSuite(testCases []TestCase, demeCount, popSize, testGenerations int, options interface{},
	suiteFunc SuiteFunc, selection []pop.SMethod,
	pairing []pop.PMethod, goal int, elites alg.IntRange, migration float64, titleSuffix string) {

	writer := SuiteWriter()

	RunSingleSuite(TestSuite{
		testCases,
		demeCount,
		popSize,
		testGenerations,
		options,
		suiteFunc,
		nil,
		nil,
		selection,
		pairing,
		goal,
		elites,
		migration,
		titleSuffix,
	}, writer)

	writer.Flush()
}

func RunSingleSuite(s TestSuite, writer *csv.Writer) {
	for _, tc := range s.testCases {
		totalGenerations := 0
		fmt.Println(tc.title)
		loops := 1
		mean := 0.0
		nextPrint := 5000
		results := []float64{}
		timings := []time.Duration{}
		fitnesses := []float64{}
		avrFitnesses := []float64{}
		for totalGenerations < s.testGenerations {

			if s.suiteInitFunc != nil {
				s.suiteInitFunc(s.initOptions)
			}

			dg := MakeDemes(
				s.demeCount,
				s.suiteFunc(s.options, s.popSize),
				s.selection,
				s.pairing,
				tc.inputs,
				tc.outputs,
				len(tc.inputs),
				s.goal,
				s.elites,
				s.migration)

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
					ind, fitness := dg.BestMember()
					fitnesses = append(fitnesses, float64(fitness))
					avrFitness := dg.AverageFitness()
					avrFitnesses = append(avrFitnesses, avrFitness)

					if loops%20 == 1 || totalGenerations > nextPrint {
						fmt.Println("Loop", loops, "Gens", totalGenerations)
						ind.Print()
						fmt.Println("Best fitness reached: ", fitness)
						fmt.Println("Generations taken: ", j+1)
						fmt.Println("Average Generations: ", float64(totalGenerations)/float64(loops))
						fmt.Println("Time taken per generation:", t2)
						fmt.Println("Average Fitness", avrFitness)
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

		timeTotal := time.Duration(0)
		for _, f := range timings {
			timeTotal += f
		}
		timeMean := timeTotal / time.Duration(len(timings))

		timeStdvTotal := 0.0
		for _, f := range timings {
			timeStdvTotal += math.Pow(float64(f-timeMean), 2)
		}
		timeStdvTotal /= float64(len(timings))
		timeStdv := time.Duration(int(math.Sqrt(timeStdvTotal)))

		fmt.Println("Average time per generation:", timeMean)
		fmt.Println("Time per generation Standard Deviation:", timeStdv)

		fitnessTotal, fitnessMean, fitnessStdv := floatSliceStatistics(fitnesses, 0)

		fmt.Println("Average best fitness:", fitnessMean)
		fmt.Println("Stdv best fitness:", fitnessStdv)

		// We should do some evening on these ranges, dropping the lowest and highest two or something.
		// This should also involve generalizing this total-mean-stdv generation for all of these sets
		// which will require mass-conversion to []float64 first.

		avrFitnessTotal, avrFitnessMean, avrFitnessStdv := floatSliceStatistics(avrFitnesses, 2)

		line := []string{
			tc.title + s.titleSuffix,
			strconv.Itoa(totalGenerations),
			floatString(mean),
			floatString(stdDev),
			timeTotal.String(),
			timeMean.String(),
			timeStdv.String(),
			floatString(fitnessTotal),
			floatString(fitnessMean),
			floatString(fitnessStdv),
			floatString(avrFitnessTotal),
			floatString(avrFitnessMean),
			floatString(avrFitnessStdv),
		}
		writer.Write(line)
	}
}

func floatString(f float64) string {
	return strconv.FormatFloat(f, 'f', 3, 64)
}

func floatSliceStatistics(fs []float64, prune int) (total float64, mean float64, stdv float64) {
	sort.Float64s(fs)

	for i := prune; i < len(fs)-prune; i++ {
		total += float64(fs[i])
	}
	mean = total / float64(len(fs)-prune)

	for i := prune; i < len(fs)-prune; i++ {
		stdv += math.Pow(float64(fs[i])-mean, 2)
	}
	stdv /= float64(len(fs) - prune)
	stdv = math.Sqrt(stdv)
	return
}
