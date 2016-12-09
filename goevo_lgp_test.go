package goevo

import (
	"fmt"
	"goevo/env"
	"goevo/lgp"
	"goevo/pairing"
	"goevo/population"
	"goevo/selection"
	"math/rand"
	"testing"
	"time"
)

func TestLGPRun(t *testing.T) {

	rand.Seed(time.Now().UTC().UnixNano())

	// Experimenting with this syntax.
	// It doesn't look very much like go right now.
	gpOpt := lgp.LGPOptions{
		MaxActionCount:  20,
		MaxStartActions: 10,
		MinStartActions: 3,

		SwapMutationChance:   0.05,
		ValueMutationChance:  0.05,
		ShrinkMutationChance: 0.05,
		ExpandMutationChance: 0.05,
		MemMutationChance:    0.05,
	}

	actions := lgp.BaseActions

	e := env.NewI(1, 0)
	mem := env.NewI(2, 0)

	in := make([][]float64, 3)
	in[0] = []float64{3.0}
	in[1] = []float64{2.0}
	in[2] = []float64{4.0}
	out := make([][]float64, 3)
	out[0] = []float64{27.0}
	out[1] = []float64{8.0}
	out[2] = []float64{64.0}

	lgp.Init(gpOpt, e, mem, lgp.PointCrossover{2},
		actions, 1.0, lgp.ComplexityFitness(lgp.Mem0Fitness, 0.1))

	lgp.AddEnvironmentAccess(1.0)

	popSize := 200
	demeCount := 5
	numGens := 10000

	members := make([][]population.Individual, demeCount)
	for j := 0; j < demeCount; j++ {
		members[j] = make([]population.Individual, popSize/demeCount)
		for i := 0; i < popSize/demeCount; i++ {
			members[j][i] = lgp.GenerateLGP(gpOpt)
		}
	}
	s := selection.DeterministicTournamentSelection{
		2,
		3,
	}

	pair := pairing.RandomPairing{}

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
		fmt.Println("Gen", i+1)
		stopEarly := dg.NextGeneration()
		if i == numGens-1 || stopEarly {
			for _, p := range dg.Demes {
				w, _ := p.Weights(1.0)
				fmt.Println(w)
				maxWeight := 0.0
				maxIndex := 0
				for i, v := range w {
					if v > maxWeight {
						maxWeight = v
						maxIndex = i
					}
				}
				fmt.Println(p.Fitnesses)
				p.Members[maxIndex].Print()
			}
			fmt.Println("Generations taken: ", i+1)
			break
		}
	}

}
