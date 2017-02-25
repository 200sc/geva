package goevo

import (
	"bitbucket.org/StephenPatrick/goevo/alg"
	"bitbucket.org/StephenPatrick/goevo/env"
	"bitbucket.org/StephenPatrick/goevo/gp"
	"bitbucket.org/StephenPatrick/goevo/lgp"
	"bitbucket.org/StephenPatrick/goevo/neural"
	"bitbucket.org/StephenPatrick/goevo/pairing"
	"bitbucket.org/StephenPatrick/goevo/pop"
	"bitbucket.org/StephenPatrick/goevo/selection"
	"testing"
)

func TestGPPow8(t *testing.T) {

	Seed()

	testCases := []TestCase{Pow8TestCase()}

	gpOpt := gp.Options{
		MaxNodeCount:         50,
		MaxStartDepth:        5,
		MaxDepth:             10,
		SwapMutationChance:   0.10,
		ShrinkMutationChance: 0.05,
	}

	gp.Init(
		gpOpt,
		env.NewI(1, 0),
		gp.PointCrossover{},
		gp.Pow8Actions,
		1.0,
		gp.ComplexityFitness(gp.Mem0Fitness, 0.01))

	gp.AddStorage(3, 1)

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
		alg.LinearIntRange{4, 6},
		0.05,
		"TGP")
}

// VSM (virtual state machine) used instead of LGP
// for regex simplicity to select GPSuite by itself
func TestVSMPow8(t *testing.T) {

	Seed()

	testCases := []TestCase{Pow8TestCase()}

	gpOpt := lgp.Options{
		MinActionCount:  2,
		MaxActionCount:  20,
		MaxStartActions: 10,
		MinStartActions: 5,

		SwapMutationChance:   0.10,
		ValueMutationChance:  0.10,
		ShrinkMutationChance: 0.10,
		ExpandMutationChance: 0.10,
		MemMutationChance:    0.10,
	}

	actions := lgp.BaseActions
	actions = append(actions, lgp.EnvActions...)

	lgp.Init(gpOpt,
		env.NewI(1, 0),
		env.NewI(2, 0),
		//lgp.UniformCrossover{0.5},
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
		alg.LinearIntRange{4, 6},
		0.05,
		"LGP")
}

func TestNNPow8(t *testing.T) {

	Seed()

	testCases := []TestCase{Pow8TestCase()}

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
