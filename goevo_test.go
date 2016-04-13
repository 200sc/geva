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
		0.40,
		0.20,
		5,
	}

	cgOpt := neural.RectifierColumnGenerationOptions{
		3,
		4,
		0.1,
	}

	nnmOpt := neural.RectifierNetworkMutationOptions{
		&wOpt,
		&cgOpt,
		0.02,
		0.00,
		0.06,
		0.00,
		0.00,
		0.33,
	}

	nngOpt := neural.RectifierNetworkGenerationOptions{
		nnmOpt,
		3,
		4,
		3,
		1,
		25,
	}

	members := make([]neural.Network, 10)
	for i := 0; i < 10; i++ {
		members[i] = nngOpt.Generate()
	}

	s := selection.TournamentSelection{
		2,
		2,
		0.8,
	}

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
		10,
		s,
		c,
		in,
		out,
		make([]int, 1),
		make([]int, 1),
	}

	p.Print()
	for i := 0; i < 10; i++ {
		fmt.Println("Gen", i)
		p = *(p.NextGeneration())
		p.Print()
	}
}
