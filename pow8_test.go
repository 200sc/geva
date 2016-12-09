package goevo

import (
	"goevo/env"
	"goevo/gp"
	"goevo/lgp"
	"goevo/pairing"
	"goevo/selection"
	"testing"
)

func TestGPSuite(t *testing.T) {

	Seed()

	testCases := make([]GPTestCase, 0)
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

	gp.AddEnvironmentAccess(1.0)

	//	gp.AddStorage(1, 1.0)

	RunSuite(
		testCases,
		5,
		200,
		100000,
		gpOpt,
		gp.GeneratePopulation,
		selection.DeterministicTournamentSelection{2, 3},
		pairing.RandomPairing{},
		5,
		1,
		0.05)
}

// VSM (virtual state machine) used instead of LGP
// for regex simplicity to select GPSuite by itself
func TestVSMSuite(t *testing.T) {

	Seed()

	testCases := make([]GPTestCase, 0)
	testCases = append(testCases, Pow8TestCase())

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

	lgp.Init(gpOpt,
		env.NewI(1, 0),
		env.NewI(2, 0),
		lgp.PointCrossover{2},
		lgp.BaseActions,
		1.0,
		lgp.ComplexityFitness(lgp.Mem0Fitness, 0.1))

	lgp.AddEnvironmentAccess(1.0)

	RunSuite(
		testCases,
		5,
		200,
		100000,
		gpOpt,
		lgp.GeneratePopulation,
		selection.DeterministicTournamentSelection{2, 3},
		pairing.RandomPairing{},
		5,
		1,
		0.05)
}
