package goevo

import (
	"fmt"
	"goevo/crossover"
	"goevo/neural"
	"goevo/pairing"
	"goevo/population"
	"goevo/selection"
	"math/rand"
	"testing"
	"time"
)

func TestPopulationRun(t *testing.T) {

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
	numGens := 100

	members := make([]neural.Network, popSize)
	for i := 0; i < popSize; i++ {
		members[i] = nngOpt.Generate()
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

	c := crossover.AverageCrossover{
		2,
	}
	// c := crossover.PointCrossover{
	// 	1,
	// }
	// c := crossover.UniformCrossover{
	// 	0.75,
	// }

	pair := pairing.RandomPairing{}

	in := make([][]float64, 4)
	in[0] = []float64{3.0, 2.0, 0.0}
	in[1] = []float64{10.0, 20.0, 10.0}
	in[2] = []float64{2.0, 100.0, 1.0}
	in[3] = []float64{0.0, 0.0, 50.0}
	out := make([][]float64, 4)
	out[0] = []float64{15.0}
	out[1] = []float64{120.0}
	out[2] = []float64{309.0}
	out[3] = []float64{150.0}

	p := population.Population{
		members,
		&nngOpt,
		popSize,
		s,
		c,
		pair,
		in,
		out,
		5,
		make([]int, popSize),
		0,
		0,
	}

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
		fmt.Println(p.Fitnesses)
		p.Members[maxIndex].Print()
	}
}
