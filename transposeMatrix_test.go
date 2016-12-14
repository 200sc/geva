package goevo

import (
	"goevo/alg"
	"goevo/env"
	"goevo/gp"
	"goevo/lgp"
	"goevo/pairing"
	"goevo/population"
	"goevo/selection"
	"testing"
)

func TestGPTransposeMatrix(t *testing.T) {
	Seed()

	testCases := make([]TestCase, 0)
	testCases = append(testCases, TransposeMatrixTestCase())

	gpOpt := gp.GPOptions{
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
		gp.BaseActions,
		1.0,
		gp.MatchMemFitness)

	gp.AddStorage(10, 1.0)

	RunSuite(
		testCases,
		5,
		200,
		100000,
		gpOpt,
		gp.GeneratePopulation,
		[]population.SelectionMethod{selection.DeterministicTournament{2, 3}},
		[]population.PairingMethod{pairing.Random{}},
		1,
		alg.LinearIntRange{4, 6},
		0.05)
}

func TestVSMTransposeMatrix(t *testing.T) {

	Seed()

	testCases := make([]TestCase, 0)
	testCases = append(testCases, TransposeMatrixTestCase())

	gpOpt := lgp.LGPOptions{
		MinActionCount:  10,
		MaxActionCount:  200,
		MaxStartActions: 40,
		MinStartActions: 20,

		SwapMutationChance:   0.15,
		ValueMutationChance:  0.15,
		ShrinkMutationChance: 0.10,
		ExpandMutationChance: 0.10,
		MemMutationChance:    0.10,
	}

	actions := lgp.BaseActions
	actions = append(actions, lgp.EnvActions...)

	lgp.Init(gpOpt,
		env.NewI(5, 0),
		env.NewI(10, 0),
		lgp.PointCrossover{3},
		actions,
		1.0,
		lgp.MatchMemFitness,
		600)

	lgp.PrintActions()

	RunSuite(
		testCases,
		75,
		3000,
		100000,
		gpOpt,
		lgp.GeneratePopulation,
		[]population.SelectionMethod{
			selection.Probabilistic{3, 2},
			selection.Probabilistic{2, 2},
			selection.DeterministicTournament{2, 3},
			selection.DeterministicTournament{3, 3},
			selection.Tournament{4, 3, 0.5},
		},
		[]population.PairingMethod{pairing.Random{}},
		1,
		alg.LinearIntRange{1, 10},
		0.10)
}
