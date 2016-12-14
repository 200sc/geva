package goevo

import (
	"fmt"
	"goevo/alg"
	"goevo/env"
	"goevo/gp"
	"goevo/pairing"
	"goevo/population"
	"goevo/selection"
	"math/rand"
	"testing"
)

const (
	actionCount = 36
	position    = 37
	direction   = 38
)

var (
	corners = map[int]bool{
		0:  true,
		5:  true,
		30: true,
		35: true,
	}
)

func forwardPosition(pos, dir int) int {
	if pos > 35 {
		return -1
	}
	switch dir {
	case 0:
		if pos < 6 {
			return -1
		}
		return pos - 6
	case 1:
		if pos%5 == 0 {
			return -1
		}
		return pos + 1
	case 2:
		if pos > 29 {
			return -1
		}
		return pos + 6
	case 3:
		if pos%6 == 0 {
			return -1
		}
		return pos - 1
	}
	return -1
}

func turn(g *gp.GP, nothing ...*gp.Node) int {
	// Increment the number of actions we've taken
	*(*g.Env)[actionCount] = *(*g.Env)[actionCount] + 1
	// 0 - N
	// 1 - E
	// 2 - S
	// 3 - W
	*(*g.Env)[direction] = (*(*g.Env)[direction] + 1) % 4
	return 1
}

func forward(g *gp.GP, nothing ...*gp.Node) int {
	// Increment the number of actions we've taken
	*(*g.Env)[actionCount] = *(*g.Env)[actionCount] + 1
	pos := *(*g.Env)[position]
	newPos := forwardPosition(pos, *(*g.Env)[direction])
	if newPos == -1 {
		// We can't walk into a wall
		return -1
	}
	v := *(*g.Env)[newPos]
	if v == 1 {
		// There's a block in the way
		// See if we can push the block
		blockPos := forwardPosition(newPos, *(*g.Env)[direction])
		if blockPos == -1 {
			// We can't push the block into a wall
			return -1
		}
		if *(*g.Env)[blockPos] == 1 {
			// We can't push the block into another block
			return -1
		}
		*(*g.Env)[newPos] = 0
		*(*g.Env)[blockPos] = 1
	}
	*(*g.Env)[position] = newPos
	return 1
}

func look(g *gp.GP, nothing ...*gp.Node) int {
	// Increment the number of actions we've taken
	//*(*g.Env)[actionCount] = *(*g.Env)[actionCount] + 1

	pos := *(*g.Env)[position]
	// -1 represents outside the map
	// 0 represents nothing
	// 1 represents a block
	newPos := forwardPosition(pos, *(*g.Env)[direction])
	if newPos == -1 {
		return newPos
	}
	return *(*g.Env)[newPos]
}

func GPTartarusFitness(g *gp.GP, inputs, outputs [][]float64) int {
	fitness := 0
	for _, envDiff := range inputs {
		g.Env = env.NewI(39, 0).New(envDiff)
		runs := 0
		*(*g.Env)[actionCount] = 0
		for runs < 200 && *(*g.Env)[actionCount] < 100 {
			gp.Eval(g.First)
			runs++
		}
		tFitness := 13

		for j, v := range *g.Env {
			if _, ok := corners[j]; ok {
				if *v == 1 {
					tFitness -= 3
				}
			} else if (j != 0 && j < 5) || (j > 30 && j < 35) ||
				j%6 == 0 || j%5 == 0 {
				if *v == 1 {
					tFitness--
				}
			}
		}
		fitness += tFitness
	}
	fitness /= 4
	return fitness
}

func RandomTartarusBoard() []float64 {
	out := make([]float64, 39)
	for i := 0; i < 4; i++ {
		v := rand.Intn(actionCount)
		for v < 6 || v > 29 || v%6 == 0 || v%5 == 0 || out[v] == 1.0 {
			v = rand.Intn(actionCount)
		}
		out[v] = 1.0
	}
	v := rand.Intn(actionCount)
	for out[v] == 1.0 {
		v = rand.Intn(actionCount)
	}
	out[position] = float64(v)
	return out
}

func PrintBoard(board env.I) {
	if len(board) < 39 {
		return
	}
	for i := 0; i < 6; i++ {
		for j := 0; j < 6; j++ {
			k := (i * 6) + j
			fmt.Print(" ", *board[k], " ")
		}
		fmt.Println("-")
	}
	fmt.Println(*board[actionCount])
	fmt.Println(*board[position])
	fmt.Println(*board[direction])
}

func TestGPTartarus(t *testing.T) {

	Seed()

	gpOpt := gp.GPOptions{
		MaxNodeCount:         100,
		MaxStartDepth:        10,
		MaxDepth:             20,
		SwapMutationChance:   0.10,
		ShrinkMutationChance: 0.05,
	}

	tartActions := []gp.Action{
		{look, "look"},
		{forward, "forward"},
		{turn, "turn"},
	}

	actions := gp.TartarusActions
	actions[0] = append(actions[0], tartActions...)

	env := env.NewI(39, 0)

	in := make([][]float64, 20)
	out := make([][]float64, 20)
	for i := 0; i < 20; i++ {
		in[i] = RandomTartarusBoard()
		out[i] = []float64{}
	}

	gp.Init(gpOpt, env, gp.PointCrossover{},
		actions, 1.0, gp.ComplexityFitness(GPTartarusFitness, 0.02))

	members := make([]population.Individual, 200)
	for j := range members {
		members[j] = gp.GenerateGP(gpOpt)
	}

	dg := MakeDemes(
		10,
		members,
		[]population.SelectionMethod{selection.Probabilistic{4, 2}},
		[]population.PairingMethod{pairing.Random{}},
		in,
		out,
		4,
		2,
		alg.LinearIntRange{1, 3},
		0.05)

	RunDemeGroup(dg, 500)
}
