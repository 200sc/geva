package goevo

import (
	"fmt"
	"goevo/crossover"
	"goevo/neural"
	"goevo/population"
	"goevo/selection"
	"testing"
)

func TestPopulationRun(t *testing.T) {
	wOpt := neural.FloatMutationOptions{
		0.50,
		0.10,
		5,
	}

	cgOpt := neural.RectifierColumnGenerationOptions{
		3,
		4,
		0.5,
	}

	nnmOpt := neural.RectifierNetworkMutationOptions{
		&wOpt,
		&cgOpt,
		0.05,
		0.00,
		0.05,
		0.00,
		0.00,
		0.30,
	}

	nngOpt := neural.RectifierNetworkGenerationOptions{
		nnmOpt,
		3,
		4,
		3,
		1,
		25,
	}

	popSize := 50

	members := make([]neural.Network, popSize)
	for i := 0; i < popSize; i++ {
		members[i] = nngOpt.Generate()
	}
	s := selection.StochasticUniversalSelection{
		2,
		false,
		1.0,
	}

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

	in := make([][]float64, 1)
	in[0] = []float64{3.0, 2.0, -1.0}
	out := make([][]float64, 1)
	out[0] = []float64{4.0}

	p := population.Population{
		members,
		&nngOpt,
		popSize,
		s,
		c,
		in,
		out,
		make([]int, 1),
		make([]int, 1),
	}

	p.Print()
	for i := 0; i < 100; i++ {
		fmt.Println("Gen", i)
		p = *(p.NextGeneration())
		fmt.Println(p.Members[0].Fitness(p.TestInputs, p.TestExpected))
		fmt.Println(p.Members[1].Fitness(p.TestInputs, p.TestExpected))
		fmt.Println(p.Members[2].Fitness(p.TestInputs, p.TestExpected))
		fmt.Println(p.Members[3].Fitness(p.TestInputs, p.TestExpected))
		fmt.Println(p.Members[4].Fitness(p.TestInputs, p.TestExpected))
	}
}
