package goevo

import (
	"goevo/alg"
	"goevo/env"
	"goevo/gp"
	"goevo/neural"
	"goevo/pairing"
	"goevo/population"
	"goevo/selection"
	"testing"
)

func TestGPRun(t *testing.T) {

	Seed()

	gpOpt := gp.GPOptions{
		MaxNodeCount:         50,
		MaxStartDepth:        5,
		MaxDepth:             10,
		SwapMutationChance:   0.10,
		ShrinkMutationChance: 0.05,
	}

	actions := gp.BaseActions

	env := env.NewI(1, 0)

	in := make([][]float64, 3)
	in[0] = []float64{3.0}
	in[1] = []float64{2.0}
	in[2] = []float64{4.0}
	out := make([][]float64, 3)
	out[0] = []float64{27.0}
	out[1] = []float64{8.0}
	out[2] = []float64{64.0}

	gp.Init(gpOpt, env, gp.PointCrossover{},
		actions, 1.0, gp.ComplexityFitness(gp.OutputFitness, 0.1))

	members := make([]population.Individual, 200)
	for j := 0; j < 200; j++ {
		members[j] = gp.GenerateGP(gpOpt)
	}

	RunDemeGroup(
		MakeDemes(
			5,
			members,
			[]population.SelectionMethod{selection.DeterministicTournament{2, 3}},
			[]population.PairingMethod{pairing.Alpha{2}},
			in,
			out,
			len(in),
			1,
			alg.LinearIntRange{1, 2},
			0.05,
		),
		10000)
}

func TestNNRun(t *testing.T) {

	Seed()

	nngOpt := neural.NetworkGenerationOptions{
		NetworkMutationOptions: neural.NetworkMutationOptions{
			WeightOptions: &neural.FloatMutationOptions{
				MutChance:     0.20,
				MutMagnitude:  0.05,
				MutRange:      20,
				ZeroOutChance: 0.01,
			},
			ColumnOptions: &neural.ColumnGenerationOptions{
				MinSize:           3,
				MaxSize:           4,
				DefaultAxonWeight: 0.5,
			},
			ActivatorOptions: &neural.ActivatorMutationOptions{
				neural.Rectifier,
				neural.Identity,
				neural.BentIdentity,
				neural.Softplus,
				neural.Softstep,
				neural.Softsign,
				neural.Sinc,
				neural.Perceptron_Threshold(0.5),
				neural.Rectifier_Exponential(1.5),
			},
			NeuronReplacementChance: 0.05,
			NeuronAdditionChance:    0.00,
			WeightSwapChance:        0.05,
			ColumnRemovalChance:     0.00,
			ColumnAdditionChance:    0.00,
			NeuronMutationChance:    0.10,
			ActivatorMutationChance: 0.01,
		},
		MinColumns:    1,
		MaxColumns:    2,
		Inputs:        3,
		Outputs:       1,
		BaseMutations: 20,
		Activator:     neural.Rectifier,
	}

	members := make([]population.Individual, 200)
	for j := range members {
		members[j] = nngOpt.Generate()
	}

	in := [][]float64{
		{3.0, 2.0, 0.0},
		{10.0, 20.0, 10.0},
		{2.0, 100.0, 1.0},
		{0.0, 0.0, 50.0},
		{10.0, 1.0, 1.0},
	}
	out := [][]float64{
		{15.0},
		{120.0},
		{309.0},
		{150.0},
		{36.0},
	}

	neural.Init(nngOpt, neural.AverageCrossover{2})

	RunDemeGroup(
		MakeDemes(
			4,
			members,
			[]population.SelectionMethod{selection.Probabilistic{3, 1.7}},
			[]population.PairingMethod{pairing.Random{}},
			in,
			out,
			len(in),
			2.0,
			alg.LinearIntRange{2, 3},
			0.1,
		),
		500)
}
