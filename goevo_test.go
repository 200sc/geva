package goevo

import (
	"fmt"
	"goevo/gp"
	"goevo/neural"
	"goevo/pairing"
	"goevo/population"
	"goevo/selection"
	"math/rand"
	"testing"
	"time"
)

func TestGPRun(t *testing.T) {

	rand.Seed(time.Now().UTC().UnixNano())

	// Experimenting with this syntax.
	// It doesn't look very much like go right now.
	gpOpt := gp.GPOptions{
		MaxNodeCount:         20,
		MaxStartDepth:        5,
		MaxDepth:             10,
		SwapMutationChance:   0.02,
		ShrinkMutationChance: 0.02,
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

	//gp.Init(gpOpt, env, gp.PointCrossover{}, actions, gp.ComplexityFitness(gp.OutputFitness, 1.0))
	gp.Init(gpOpt, env, gp.PointCrossover{}, actions, gp.OutputFitness)
	gp.AddEnvironmentAccess()

	popSize := 100
	demeCount := 5
	numGens := 200

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
			Elites:       2,
			Fitnesses:    make([]int, popSize/demeCount),
		}
	}
	dg := population.DemeGroup{
		Demes:           demes,
		MigrationChance: 0.05,
	}

	for i := 0; i < numGens; i++ {
		fmt.Println("Gen", i)
		dg.NextGeneration()
		if i == numGens-1 {
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
				p.Members[0].Print()
				p.Members[1].Print()
				p.Members[3].Print()
				p.Members[maxIndex].Print()
			}
		}
	}

}

func TestNNRun(t *testing.T) {

	rand.Seed(time.Now().UTC().UnixNano())

	wOpt := neural.FloatMutationOptions{
		0.20,
		0.05,
		20,
		0.01,
	}

	cgOpt := neural.ColumnGenerationOptions{
		3,
		4,
		0.5,
	}

	actOpt := neural.ActivatorMutationOptions{
		neural.Rectifier,
		neural.Identity,
		neural.BentIdentity,
		neural.Softplus,
		neural.Softstep,
		neural.Softsign,
		neural.Sinc,
		neural.Perceptron_Threshold(0.5),
		neural.Rectifier_Exponential(1.5),
	}

	nnmOpt := neural.NetworkMutationOptions{
		&wOpt,
		&cgOpt,
		&actOpt,
		0.05,
		0.00,
		0.05,
		0.00,
		0.00,
		0.10,
		0.01,
	}

	nngOpt := neural.NetworkGenerationOptions{
		nnmOpt,
		1,
		2,
		3,
		1,
		20,
		neural.Rectifier,
	}

	popSize := 200
	demeCount := 4
	numGens := 500

	members := make([][]population.Individual, demeCount)
	for j := 0; j < demeCount; j++ {
		members[j] = make([]population.Individual, popSize/demeCount)
		for i := 0; i < popSize/demeCount; i++ {
			members[j][i] = nngOpt.Generate()
		}
	}
	s := selection.ProbabilisticSelection{
		3,
		1.7,
	}
	// s := selection.StochasticUniversalSelection{
	// 	2,
	// 	false,
	// 	1.0,
	// }
	// s := selection.TournamentSelection{
	// 	2,
	// 	2,
	// 	1.0,
	// }
	// s := selection.GreedySelection{
	// 	2,
	// }

	pair := pairing.RandomPairing{}

	in := make([][]float64, 5)
	in[0] = []float64{3.0, 2.0, 0.0}
	in[1] = []float64{10.0, 20.0, 10.0}
	in[2] = []float64{2.0, 100.0, 1.0}
	in[3] = []float64{0.0, 0.0, 50.0}
	in[4] = []float64{10.0, 1.0, 1.0}
	out := make([][]float64, 5)
	out[0] = []float64{15.0}
	out[1] = []float64{120.0}
	out[2] = []float64{309.0}
	out[3] = []float64{150.0}
	out[4] = []float64{36.0}

	demes := make([]population.Population, demeCount)
	for i := 0; i < demeCount; i++ {
		demes[i] = population.Population{
			Members:      members[i],
			Size:         popSize / demeCount,
			Selection:    s,
			Pairing:      pair,
			FitnessTests: 5,
			TestInputs:   in,
			TestExpected: out,
			Elites:       2,
			Fitnesses:    make([]int, popSize/demeCount),
		}
	}
	dg := population.DemeGroup{
		Demes:           demes,
		MigrationChance: 0.1,
	}

	neural.Init(nngOpt, neural.AverageCrossover{2})

	for i := 0; i < numGens; i++ {
		fmt.Println("Gen", i)
		dg.NextGeneration()
		if i == numGens-1 {
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
		}
	}
}
