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

func TestGPPow8(t *testing.T) {

	Seed()

	testCases := make([]TestCase, 0)
	testCases = append(testCases, Pow8TestCase())

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
		gp.ComplexityFitness(gp.OutputFitness, 0.05))

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

// VSM (virtual state machine) used instead of LGP
// for regex simplicity to select GPSuite by itself
func TestVSMPow8(t *testing.T) {

	Seed()

	testCases := make([]TestCase, 0)
	testCases = append(testCases, Pow8TestCase())

	gpOpt := lgp.LGPOptions{
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
		[]population.SelectionMethod{selection.DeterministicTournament{2, 3}},
		[]population.PairingMethod{pairing.Random{}},
		1,
		alg.LinearIntRange{4, 6},
		0.05)
}
