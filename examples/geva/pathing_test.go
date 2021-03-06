package geva

import (
	"math"
	"testing"

	"github.com/200sc/go-dist/intrange"

	"github.com/200sc/geva/env"
	"github.com/200sc/geva/gp"
	"github.com/200sc/geva/lgp"
	"github.com/200sc/geva/pairing"
	"github.com/200sc/geva/pop"
	"github.com/200sc/geva/selection"
)

func lookLGPvar(g *lgp.LGP, xs ...int) {
	g.SetReg(xs[0], look(g.Env, g.RegVal(xs[1])))
}

func lookGPvar(g *gp.GP, xs ...*gp.Node) int {
	return look(g.Env, gp.Eval(xs[0]))
}

func PathingFitnessGP(g *gp.GP, inputs, outputs [][]float64) int {
	fitness := 0
	for _, envDiff := range inputs {
		g.Env = env.NewI(39, 0).New(envDiff)
		runs := 0
		*(*g.Env)[actionCount] = 0
		for runs < 200 && *(*g.Env)[actionCount] < 100 && *(*g.Env)[position] != 35 {
			gp.Eval(g.First)
			runs++
		}
		pos := *(*g.Env)[position]
		x1 := pos % 6
		y1 := pos / 6
		fitness += int(intDistance(x1, y1, 5, 5) * 100)
		fitness += *(*g.Env)[actionCount]
	}
	fitness /= len(inputs)
	return fitness
}

func intDistance(x1, y1, x2, y2 int) float64 {
	return math.Sqrt(
		math.Pow(float64(x1-x2), 2) +
			math.Pow(float64(y1-y2), 2))
}

func PathingFitnessLGP(g *lgp.LGP, inputs, outputs [][]float64) int {
	fitness := 0
	for _, envDiff := range inputs {
		g.Env = env.NewI(39, 0).New(envDiff)

		runs := 0
		*(*g.Env)[actionCount] = 0
		for runs < 200 && *(*g.Env)[actionCount] < 100 && *(*g.Env)[position] != 35 {
			g.Run()
			runs++
		}
		pos := *(*g.Env)[position]
		x1 := pos % 6
		y1 := pos / 6
		fitness += int(intDistance(x1, y1, 5, 5) * 100)
		fitness += *(*g.Env)[actionCount]
	}
	fitness /= len(inputs)
	return fitness
}

func TestGPPathing(t *testing.T) {

	Seed()

	gpOpt := gp.Options{
		MaxNodeCount:         100,
		MaxStartDepth:        10,
		MaxDepth:             20,
		SwapMutationChance:   0.10,
		ShrinkMutationChance: 0.05,
	}

	tartActions := []gp.Action{
		{forwardGP, "forward"},
		{turnGP, "turn"},
	}

	actions := gp.TartarusActions
	actions[0] = append(actions[0], tartActions...)
	actions[1] = append(actions[1], gp.Action{lookGPvar, "look"})

	env := env.NewI(39, 0)

	in := make([][]float64, 20)
	out := make([][]float64, 20)
	for i := 0; i < 20; i++ {
		in[i] = RandomTartarusBoard()
		out[i] = []float64{}
	}

	gp.Init(gpOpt, env, gp.PointCrossover{},
		actions, 1.0, PathingFitnessGP)

	testCases := []TestCase{{in, out, "Pathing"}}

	RunSuite(
		testCases,
		5,
		200,
		100000,
		gpOpt,
		gp.GeneratePopulation,
		[]pop.SMethod{selection.Probabilistic{4, 2}},
		[]pop.PMethod{pairing.Random{}},
		30,
		intrange.NewLinear(1, 6),
		0.05,
		"TGP")
}

func TestLGPPathing(t *testing.T) {

	Seed()

	gpOpt := lgp.Options{
		MinActionCount:  5,
		MaxActionCount:  100,
		MaxStartActions: 50,
		MinStartActions: 20,

		SwapMutationChance:   0.10,
		ValueMutationChance:  0.10,
		ShrinkMutationChance: 0.10,
		ExpandMutationChance: 0.10,
		MemMutationChance:    0.10,
	}

	tartActions := []lgp.Action{
		{lookLGPvar, "look", 2},
		{forwardLGP, "forward", 1},
		{turnLGP, "turn", 1},
	}

	actions := lgp.TartarusActions
	actions = append(actions, tartActions...)

	in := make([][]float64, 20)
	out := make([][]float64, 20)
	for i := 0; i < 20; i++ {
		in[i] = RandomTartarusBoard()
		in[i][position] = 0
		out[i] = []float64{}
	}

	testCases := []TestCase{{in, out, "Pathing"}}

	lgp.Init(
		gpOpt,
		env.NewI(39, 0),
		env.NewI(5, 0),
		lgp.PointCrossover{2},
		actions,
		1.0,
		PathingFitnessLGP,
		200)

	RunSuite(
		testCases,
		5,
		200,
		100000,
		gpOpt,
		lgp.GeneratePopulation,
		[]pop.SMethod{selection.Probabilistic{4, 2}},
		[]pop.PMethod{pairing.Random{}},
		30,
		intrange.NewLinear(1, 6),
		0.05,
		"LGP")
}
