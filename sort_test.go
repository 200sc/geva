package goevo

import (
	"fmt"
	"goevo/alg"
	"goevo/env"
	"goevo/gp"
	"goevo/lgp"
	"goevo/pairing"
	"goevo/pop"
	"goevo/selection"
	"testing"
)

// Could be improved with more specific actions
func TestGPSortList(t *testing.T) {
	Seed()

	testCases := make([]TestCase, 0)
	testCases = append(testCases, SortListTestCase())

	gpOpt := gp.Options{
		MaxNodeCount:         250,
		MaxStartDepth:        6,
		MaxDepth:             12,
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

	gp.AddStorage(20, 1.0)

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
		0.05)
}

func TestVSMSortList(t *testing.T) {

	Seed()

	testCases := make([]TestCase, 0)
	testCases = append(testCases, SortListTestCase())

	fmt.Println(testCases)

	gpOpt := lgp.Options{
		MinActionCount:  10,
		MaxActionCount:  200,
		MaxStartActions: 40,
		MinStartActions: 20,

		SwapMutationChance:   0.15,
		ValueMutationChance:  0.15,
		ShrinkMutationChance: 0.10,
		ExpandMutationChance: 0.10,
		MemMutationChance:    0.00,
	}

	actions := lgp.SortActions

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
		[]pop.SMethod{
			selection.Probabilistic{3, 2},
			selection.Probabilistic{2, 2},
			selection.DeterministicTournament{2, 3},
			selection.DeterministicTournament{3, 3},
			selection.Tournament{4, 3, 0.5},
		},
		[]pop.PMethod{pairing.Random{}},
		1,
		alg.LinearIntRange{1, 10},
		0.10)
}
