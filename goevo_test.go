package goevo

import (
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
		2,
		16,
		0.1,
	}

	nnmOpt := neural.RectifierNetworkMutationOptions{
		&wOpt,
		&cgOpt,
		0.02,
		0.06,
		0.06,
		0.01,
		0.01,
		0.33,
	}

	nngOpt := neural.RectifierNetworkGenerationOptions{
		nnmOpt,
		5,
		20,
		3,
		4,
		50,
	}

	members := make([]neural.Network, 10)
	for i := 0; i < 10; i++ {
		members[i] = neural.GenerateRectifierNetwork(&nngOpt)
	}

	s := selection.TournamentSelection{
		4,
		2,
		0.9,
	}

	c := crossover.PointCrossover{
		2,
	}

	p := population.Population{
		members,
		&nngOpt,
		10,
		s,
		c,
		make([][]float64, 1),
		make([][]float64, 1),
		make([]int, 1),
		make([]int, 1),
	}
}
