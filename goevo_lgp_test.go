package goevo

import (
	"goevo/env"
	"goevo/lgp"
	"goevo/pairing"
	"goevo/selection"
	"testing"
)

func TestLGPRun(t *testing.T) {

	Seed()

	gpOpt := lgp.LGPOptions{
		MinActionCount:  2,
		MaxActionCount:  20,
		MaxStartActions: 10,
		MinStartActions: 3,

		SwapMutationChance:   0.05,
		ValueMutationChance:  0.05,
		ShrinkMutationChance: 0.05,
		ExpandMutationChance: 0.05,
		MemMutationChance:    0.05,
	}

	in := make([][]float64, 3)
	in[0] = []float64{3.0}
	in[1] = []float64{2.0}
	in[2] = []float64{4.0}
	out := make([][]float64, 3)
	out[0] = []float64{27.0}
	out[1] = []float64{8.0}
	out[2] = []float64{64.0}

	lgp.Init(gpOpt,
		env.NewI(1, 0),
		env.NewI(2, 0),
		lgp.PointCrossover{2},
		lgp.BaseActions,
		1.0,
		lgp.ComplexityFitness(lgp.Mem0Fitness, 0.1))

	lgp.AddEnvironmentAccess(1.0)

	dg := MakeDemes(
		5,
		lgp.GeneratePopulation(gpOpt, 200),
		selection.DeterministicTournamentSelection{2, 3},
		pairing.RandomPairing{},
		in,
		out,
		3,
		1,
		1,
		0.05)

	RunDemeGroup(dg, 10000)
}