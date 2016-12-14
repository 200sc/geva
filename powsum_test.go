package goevo

import (
	"fmt"
	"goevo/alg"
	"goevo/env"
	"goevo/gp"
	"goevo/lgp"
	"goevo/pairing"
	"goevo/population"
	"goevo/selection"
	"testing"
)

func TestGPPowSum(t *testing.T) {

	Seed()

	testCases := make([]TestCase, 0)
	testCases = append(testCases, PowSumTestCase())
	fmt.Println(testCases)

	gpOpt := gp.GPOptions{
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
		[]population.SelectionMethod{selection.DeterministicTournament{2, 3}},
		[]population.PairingMethod{pairing.Random{}},
		1,
		alg.LinearIntRange{1, 10},
		0.05)
}

func TestVSMPowSum(t *testing.T) {

	Seed()

	testCases := make([]TestCase, 0)
	testCases = append(testCases, PowSumTestCase())

	gpOpt := lgp.LGPOptions{
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
		[]population.SelectionMethod{selection.DeterministicTournament{2, 3}},
		[]population.PairingMethod{pairing.Random{}},
		1,
		alg.LinearIntRange{1, 2},
		0.05)
}
