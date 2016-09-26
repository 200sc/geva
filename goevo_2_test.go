package goevo

import (
	"fmt"
	"goevo/gp"
	"goevo/pairing"
	"goevo/population"
	"goevo/selection"
	"math"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

type GPTestCase struct {
	inputs  [][]float64
	outputs [][]float64
	title   string
}

func TestGPSuite(t *testing.T) {

	testCases := make([]GPTestCase, 0)
	in := [][]float64{
		{1.0},
		{2.0},
		{3.0},
		{4.0},
	}
	// Pow 1 and Pow 2 have an average
	// gen count of 1, so they aren't
	// included here.
	for i := 3; i < 13; i++ {
		out := [][]float64{
			{math.Pow(1.0, float64(i))},
			{math.Pow(2.0, float64(i))},
			{math.Pow(3.0, float64(i))},
			{math.Pow(4.0, float64(i))},
		}
		title := "Pow" + strconv.Itoa(i)
		testCases = append(testCases, GPTestCase{
			in,
			out,
			title,
		})
	}
	// Addition tests.
	// Summary:
	// From 1-9, a few hundred generations on average.
	// From 10 onward, usually doesn't get the answer.
	// for i := 1; i <= 20; i++ {
	// 	out := [][]float64{
	// 		{1.0 + float64(i)},
	// 		{2.0 + float64(i)},
	// 		{3.0 + float64(i)},
	// 		{4.0 + float64(i)},
	// 	}
	// 	title := "Add" + strconv.Itoa(i)
	// 	testCases = append(testCases, GPTestCase{
	// 		in,
	// 		out,
	// 		title,
	// 	})
	// }

	testGenerations := 20000
	rand.Seed(time.Now().UTC().UnixNano())
	for _, tc := range testCases {
		totalGenerations := 0
		fmt.Println(tc.title)
		loops := 0
		for totalGenerations < testGenerations {

			// Experimenting with this syntax.
			// It doesn't look very much like go right now.
			gpOpt := gp.GPOptions{
				MaxNodeCount:         50,
				MaxStartDepth:        5,
				MaxDepth:             10,
				SwapMutationChance:   0.10,
				ShrinkMutationChance: 0.05,
			}

			actions := gp.BaseActions

			val := 0
			env := gp.Environment{&val}

			gp.Init(gpOpt, env, gp.PointCrossover{}, actions, 1.0,
				gp.ComplexityFitness(gp.OutputFitness, 0.05))
			//gp.Init(gpOpt, env, gp.PointCrossover{}, actions, 1.0, gp.OutputFitness)
			gp.AddEnvironmentAccess(1.0)
			gp.AddStorage(1, 1.0)

			popSize := 200
			demeCount := 5
			numGens := 10000

			members := make([][]population.Individual, demeCount)
			for j := 0; j < demeCount; j++ {
				members[j] = make([]population.Individual, popSize/demeCount)
				for i := 0; i < popSize/demeCount; i++ {
					members[j][i] = gp.GenerateGP(gpOpt)
				}
			}
			s := selection.DeterministicTournamentSelection{
				2,
				3,
			}

			pair := pairing.RandomPairing{}
			// Alpha pairing{2} doubles the expected generations to reach 1 fitness
			//pair := pairing.AlphaPairing{2}

			demes := make([]population.Population, demeCount)
			for i := 0; i < demeCount; i++ {
				demes[i] = population.Population{
					Members:      members[i],
					Size:         popSize / demeCount,
					Selection:    s,
					Pairing:      pair,
					FitnessTests: len(in),
					TestInputs:   tc.inputs,
					TestExpected: tc.outputs,
					Elites:       5,
					Fitnesses:    make([]int, popSize/demeCount),
					GoalFitness:  1,
				}
			}
			dg := population.DemeGroup{
				Demes:           demes,
				MigrationChance: 0.05,
			}

			for j := 0; j < numGens; j++ {
				stopEarly := dg.NextGeneration()
				if j == numGens-1 || stopEarly {
					totalGenerations += j + 1
					//if loops%200 == 0 {
					//fmt.Println("Loop", loops, "Gens", totalGenerations)
					//ind, _ := dg.BestMember()
					//ind.Print()
					// 	fmt.Println("Generations taken: ", j+1)
					//}
					break
				}
			}
			loops += 1
		}
		fmt.Println("Average Generations: ", float64(totalGenerations)/float64(loops))
	}
}
