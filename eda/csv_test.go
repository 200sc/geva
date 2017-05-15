package eda

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"bitbucket.org/StephenPatrick/goevo/eda/fitness"
	"bitbucket.org/StephenPatrick/goevo/mut"

	"github.com/200sc/go-compgeo/printutil"
	"github.com/200sc/go-dist/floatrange"
)

func TestCSVOut(t *testing.T) {
	file, err := os.Create("test.csv")
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	csvFile := csv.NewWriter(file)
	headers := []string{
		"TestName",
		"IterationsTaken",
		"BestFitness",
		"IterationOfBestFitness",
		"FitnessEvaluations",
		"FitnessEvaluationsAtBest",
		"TimeTaken",
		"Samples",
		"LearningSamples",
		"BaseValue",
		"LearningRate",
		"MutationRate",
	}
	csvFile.Write(headers)
	csvFile.Flush()
	csvReport := func(m Model) {
		bm := m.BaseModel()
		strs := make([]string, 12)
		strs[0] = bm.name
		strs[1] = strconv.Itoa(*bm.iterations)
		strs[2] = strconv.Itoa(*bm.bestFitness)
		strs[3] = strconv.Itoa(*bm.bestIteration)
		strs[4] = strconv.Itoa(*bm.fitnessEvals)
		strs[5] = strconv.Itoa(*bm.bestFitnessEvals)
		strs[6] = strconv.FormatInt(time.Since(bm.startTime).Nanoseconds(), 10)
		strs[7] = strconv.Itoa(bm.samples)
		strs[8] = strconv.Itoa(bm.learningSamples)
		strs[9] = printutil.Stringf64(bm.baseValue)
		strs[10] = printutil.Stringf64(bm.learningRate)
		strs[11] = printutil.Stringf64(bm.mutationRate)
		csvFile.Write(strs)
		csvFile.Flush()
		fmt.Println("Wrote", bm.name)
	}

	BenchTest = And(BenchTest, ReportFunc(csvReport))

	Tests := []Option{
		And(Length(300), FitnessFunc(fitness.OnemaxABS)),
		And(Length(300), FitnessFunc(fitness.OnemaxChance)),
		And(Length(300), FitnessFunc(fitness.AlternatingABS)),
		And(Length(100), FitnessFunc(fitness.FourPeaks(10))),
		And(Length(100), FitnessFunc(fitness.SixPeaks(10))),
		And(Length(100), FitnessFunc(fitness.Quadratic)),
		And(Length(99), FitnessFunc(fitness.TrapABS(3))),
		And(Length(100), FitnessFunc(fitness.TrapABS(5))),
	}

	TestNames := []string{
		"OnemaxABS",
		"OnemaxChance",
		"Alternating",
		"FourPeaks",
		"SixPeaks",
		"Quadratic",
		"Trap3",
		"Trap5",
	}

	ModelNames := []string{
		"PBIL",
		"CGA",
		"UMDA",
		"SHCLVND",
		"MIMIC",
		"BMDA",
		"ECGA",
		"BOA",
	}

	Models := []EDA{
		PBILModel,
		CGAModel,
		UMDAModel,
		SHCLVNDModel,
		MIMICModel,
		BMDAModel,
		ECGAModel,
		BOAModel,
	}

	ModelOpts := []Option{
		And(LearningRate(0.2), MutationRate(0.03)),
		And(LearningRate(0.1), MutationRate(0.03)),
		MutationRate(0.03),
		And(LearningRate(0.05),
			LMutator(
				mut.And(
					mut.Scale(0.997),
					EnforceRange(floatrange.NewLinear(0.001, 1.0))),
			)),
		LearningRate(0.07),
		MutationRate(.25),
		And(LearningRate(0.2), MutationRate(0.01)),
		And(LearningRate(0.1), MutationRate(0.03)),
	}

	// run forever
	for i := 0; i < 100; i++ {
		for k, t := range Tests {
			for j, m := range Models {
				mOpts := ModelOpts[j]
				Loop(m, BenchTest, t, mOpts, Name(TestNames[k]+"_"+ModelNames[j]))
			}
		}
	}
}
