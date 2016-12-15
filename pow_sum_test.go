package goevo

import (
	"fmt"
	"goevo/alg"
	"goevo/env"
	"goevo/gp"
	"goevo/lgp"
	"goevo/neural"
	"goevo/pairing"
	"goevo/pop"
	"goevo/selection"
	"testing"
)

func TestGPPowSum(t *testing.T) {

	Seed()

	testCases := make([]TestCase, 0)
	testCases = append(testCases, PowSumTestCase())
	fmt.Println(testCases)

	gpOpt := gp.Options{
		MaxNodeCount:         50,
		MaxStartDepth:        5,
		MaxDepth:             10,
		SwapMutationChance:   0.10,
		ShrinkMutationChance: 0.05,
	}

	gp.Init(
		gpOpt,
		env.NewI(2, 0),
		gp.PointCrossover{},
		gp.BaseActions,
		1.0,
		gp.OutputFitness)

	RunSuite(
		testCases,
		5,
		200,
		100000,
		gpOpt,
		gp.GeneratePopulation,
		[]pop.SMethod{selection.DeterministicTournament{2, 3}},
		[]pop.PMethod{pairing.Random{}},
		1,
		alg.LinearIntRange{1, 10},
		0.05,
		"TGP")
}

func TestVSMPowSum(t *testing.T) {

	Seed()

	testCases := make([]TestCase, 0)
	testCases = append(testCases, PowSumTestCase())

	gpOpt := lgp.Options{
		MinActionCount:  2,
		MaxActionCount:  20,
		MaxStartActions: 10,
		MinStartActions: 3,

		SwapMutationChance:   0.20,
		ValueMutationChance:  0.20,
		ShrinkMutationChance: 0.20,
		ExpandMutationChance: 0.20,
		MemMutationChance:    0.20,
	}

	actions := lgp.PowSumActions

	lgp.Init(gpOpt,
		env.NewI(1, 0),
		env.NewI(8, 0),
		lgp.PointCrossover{2},
		actions,
		1.0,
		lgp.Mem0Fitness,
		300)

	RunSuite(
		testCases,
		5,
		200,
		100000,
		gpOpt,
		lgp.GeneratePopulation,
		[]pop.SMethod{selection.DeterministicTournament{2, 3}},
		[]pop.PMethod{pairing.Random{}},
		1,
		alg.LinearIntRange{1, 2},
		0.05,
		"LGP")
}

func TestNNPowSum(t *testing.T) {

	Seed()

	testCases := []TestCase{PowSumTestCase()}

	nngOpt := neural.NetworkGenerationOptions{
		NetworkMutationOptions: neural.NetworkMutationOptions{
			WeightOptions: neural.FloatMutationOptions{
				MutChance:     0.20,
				MutMagnitude:  2.0,
				MutRange:      60,
				ZeroOutChance: 0.01,
			},
			ColumnOptions: neural.ColumnGenerationOptions{
				MinSize:           3,
				MaxSize:           4,
				DefaultAxonWeight: 0.5,
			},
			ActivatorOptions:        neural.AllActivators,
			NeuronReplacementChance: 0.05,
			NeuronAdditionChance:    0.00,
			WeightSwapChance:        0.05,
			ColumnRemovalChance:     0.00,
			ColumnAdditionChance:    0.00,
			NeuronMutationChance:    0.10,
			ActivatorMutationChance: 0.01,
		},
		MinColumns:    3,
		MaxColumns:    4,
		MaxInputs:     2,
		MaxOutputs:    1,
		BaseMutations: 20,
	}

	neural.Init(
		nngOpt,
		neural.AverageCrossover{2},
		neural.AbsFitness,
	)

	RunSuite(
		testCases,
		4,
		200,
		100000,
		nngOpt,
		neural.GeneratePopulation,
		[]pop.SMethod{selection.DeterministicTournament{3, 3}},
		[]pop.PMethod{pairing.Random{}},
		2.0,
		alg.LinearIntRange{1, 4},
		0.1,
		"ENN",
	)
}
