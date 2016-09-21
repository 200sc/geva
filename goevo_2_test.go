package goevo

import (
	"fmt"
	"goevo/gp"
	"goevo/pairing"
	"goevo/population"
	"goevo/selection"
	"math/rand"
	"testing"
	"time"
)

func TestGPAverageGenerations(t *testing.T) {
	totalGenerations := 0
	loops := 200
	rand.Seed(time.Now().UTC().UnixNano())
	for i := 0; i < loops; i++ {

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

		in := make([][]float64, 3)
		in[0] = []float64{3.0}
		in[1] = []float64{2.0}
		in[2] = []float64{4.0}
		out := make([][]float64, 3)
		out[0] = []float64{27.0}
		out[1] = []float64{8.0}
		out[2] = []float64{64.0}

		gp.Init(gpOpt, env, gp.PointCrossover{}, actions, gp.ComplexityFitness(gp.OutputFitness, 0.1))
		//gp.Init(gpOpt, env, gp.PointCrossover{}, actions, gp.OutputFitness)
		gp.AddEnvironmentAccess()

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

		pair := pairing.AlphaPairing{2}

		demes := make([]population.Population, demeCount)
		for i := 0; i < demeCount; i++ {
			demes[i] = population.Population{
				Members:      members[i],
				Size:         popSize / demeCount,
				Selection:    s,
				Pairing:      pair,
				FitnessTests: 3,
				TestInputs:   in,
				TestExpected: out,
				Elites:       1,
				Fitnesses:    make([]int, popSize/demeCount),
				GoalFitness:  1,
			}
		}
		dg := population.DemeGroup{
			Demes:           demes,
			MigrationChance: 0.05,
		}

		for i := 0; i < numGens; i++ {
			stopEarly := dg.NextGeneration()
			if i == numGens-1 || stopEarly {
				totalGenerations += i + 1
				break
			}
		}
	}
	fmt.Println("Average Generations: ", float64(totalGenerations)/float64(loops))
}
