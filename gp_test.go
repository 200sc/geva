package goevo

import (
	"bitbucket.org/StephenPatrick/goevo/alg"
	"bitbucket.org/StephenPatrick/goevo/env"
	"bitbucket.org/StephenPatrick/goevo/gp"
	"bitbucket.org/StephenPatrick/goevo/pairing"
	"bitbucket.org/StephenPatrick/goevo/pop"
	"bitbucket.org/StephenPatrick/goevo/selection"
	"math"
	"testing"
)

func TestGPRun(t *testing.T) {

	Seed()

	gpOpt := gp.Options{
		MaxNodeCount:         50,
		MaxStartDepth:        5,
		MaxDepth:             10,
		SwapMutationChance:   0.10,
		ShrinkMutationChance: 0.05,
	}

	actions := gp.BaseActions

	env := env.NewI(1, 0)

	in := [][]float64{
		{3.0},
		{2.0},
		{4.0},
	}
	out := make([][]float64, 3)
	for i, f := range in {
		out[i] = []float64{math.Pow(f[0], 3)}
	}

	gp.Init(gpOpt, env, gp.PointCrossover{},
		actions, 1.0, gp.ComplexityFitness(gp.OutputFitness, 0.1))

	members := make([]pop.Individual, 200)
	for j := 0; j < 200; j++ {
		members[j] = gp.GenerateGP(gpOpt)
	}

	RunDemeGroup(
		MakeDemes(
			5,
			members,
			[]pop.SMethod{selection.DeterministicTournament{2, 3}},
			[]pop.PMethod{pairing.Alpha{2}},
			in,
			out,
			len(in),
			1,
			alg.LinearIntRange{1, 2},
			0.05,
		),
		10000)
}
