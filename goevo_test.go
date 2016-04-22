package goevo

import (
	"fmt"
	"goevo/crossover"
	"goevo/neural"
	"goevo/population"
	"goevo/selection"
	"math"
	"math/rand"
	"testing"
	"time"
)

func TestPopulationRun(t *testing.T) {

	rand.Seed(time.Now().UTC().UnixNano())

	wOpt := neural.FloatMutationOptions{
		0.20,
		0.05,
		5,
	}

	cgOpt := neural.ModularColumnGenerationOptions{
		3,
		4,
		0.5,
	}

	nnmOpt := neural.ModularNetworkMutationOptions{
		&wOpt,
		&cgOpt,
		0.05,
		0.00,
		0.05,
		0.00,
		0.00,
		0.30,
	}

	nngOpt := neural.ModularNetworkGenerationOptions{
		nnmOpt,
		1,
		2,
		3,
		1,
		25,
		// Perceptron example
		// func(x float64) float64 {
		// 	if x > 0.5 {
		// 		return 1
		// 	}
		// 	return 0.0
		// },
		// Rectifier example
		// (far better suited for our test problem)
		func(x float64) float64 {
			return math.Max(x, 0.0)
		},
	}

	popSize := 100
	numGens := 300

	members := make([]neural.ModularNetwork, popSize)
	for i := 0; i < popSize; i++ {
		members[i] = nngOpt.Generate()
	}
	s := selection.ProbabilisticSelection{
		2,
		1.4,
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

	c := crossover.PointCrossover{
		1,
	}
	// c := crossover.UniformCrossover{
	// 	0.5,
	// }

	in := make([][]float64, 3)
	in[0] = []float64{3.0, 2.0, 0.0}
	in[1] = []float64{10.0, 20.0, 10.0}
	in[2] = []float64{2.0, 100.0, 1.0}
	out := make([][]float64, 3)
	out[0] = []float64{5.0}
	out[1] = []float64{40.0}
	out[2] = []float64{103.0}

	p := population.Population{
		members,
		&nngOpt,
		popSize,
		s,
		c,
		in,
		out,
	}

	p.Print()
	for i := 0; i < numGens; i++ {
		fmt.Println("Gen", i)
		p = *(p.NextGeneration())
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
		fmt.Println(p.Members[maxIndex].Fitness(p.TestInputs, p.TestExpected))
		p.Members[maxIndex].Print()
	}
}
