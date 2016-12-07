package goevo

import (
	"fmt"
	"goevo/env"
	"goevo/gp"
	"goevo/pairing"
	"goevo/population"
	"goevo/selection"
	"math/rand"
	"testing"
	"time"
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
	*(*g.Env)[36] = *(*g.Env)[36] + 1
	// 0 - N
	// 1 - E
	// 2 - S
	// 3 - W
	*(*g.Env)[38] = (*(*g.Env)[38] + 1) % 4
	return 1
}

func forward(g *gp.GP, nothing ...*gp.Node) int {
	// Increment the number of actions we've taken
	*(*g.Env)[36] = *(*g.Env)[36] + 1
	pos := *(*g.Env)[37]
	newPos := forwardPosition(pos, *(*g.Env)[38])
	if newPos == -1 {
		// We can't walk into a wall
		return -1
	}
	v := *(*g.Env)[newPos]
	if v == 1 {
		// There's a block in the way
		// See if we can push the block
		blockPos := forwardPosition(newPos, *(*g.Env)[38])
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
	*(*g.Env)[37] = newPos
	return 1
}

func look(g *gp.GP, nothing ...*gp.Node) int {
	// Increment the number of actions we've taken
	//*(*g.Env)[36] = *(*g.Env)[36] + 1

	pos := *(*g.Env)[37]
	// -1 represents outside the map
	// 0 represents nothing
	// 1 represents a block
	newPos := forwardPosition(pos, *(*g.Env)[38])
	if newPos == -1 {
		return newPos
	}
	return *(*g.Env)[newPos]
}

func TartarusFitness(g *gp.GP, inputs, outputs [][]float64) int {
	fitness := 0
	for _, envDiff := range inputs {
		g.Env = env.NewI(39, 0).New(envDiff)
		runs := 0
		*(*g.Env)[36] = 0
		for runs < 200 && *(*g.Env)[36] < 100 {
			gp.Eval(g.First)
			runs++
		}
		t_fitness := 13

		for j, v := range *g.Env {
			if _, ok := corners[j]; ok {
				if *v == 1 {
					t_fitness -= 3
				}
			} else if (j != 0 && j < 5) || (j > 30 && j < 35) ||
				j%6 == 0 || j%5 == 0 {
				if *v == 1 {
					t_fitness -= 1
				}
			}
		}
		fitness += t_fitness
	}
	fitness /= 4
	return fitness
}

func RandomTartarusBoard() []float64 {
	out := make([]float64, 39)
	for i := 0; i < 4; i++ {
		v := rand.Intn(36)
		for v < 6 || v > 29 || v%6 == 0 || v%5 == 0 || out[v] == 1.0 {
			v = rand.Intn(36)
		}
		out[v] = 1.0
	}
	v := rand.Intn(36)
	for out[v] == 1.0 {
		v = rand.Intn(36)
	}
	out[37] = float64(v)
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
	fmt.Println(*board[36])
	fmt.Println(*board[37])
	fmt.Println(*board[38])
}

func TestGPTartarus(t *testing.T) {

	rand.Seed(time.Now().UTC().UnixNano())

	// Experimenting with this syntax.
	// It doesn't look very much like go right now.
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

	actions := gp.BaseActions
	actions[0] = append(actions[0], tartActions...)

	env := env.NewI(39, 0)

	in := make([][]float64, 20)
	out := make([][]float64, 20)
	for i := 0; i < 20; i++ {
		in[i] = RandomTartarusBoard()
		// Tartarus doesn't have an out comparison
		out[i] = []float64{}
	}

	gp.Init(gpOpt, env, gp.PointCrossover{},
		actions, 1.0, gp.ComplexityFitness(TartarusFitness, 0.02))
	//gp.Init(gpOpt, env, gp.PointCrossover{}, actions, gp.OutputFitness)
	//gp.AddEnvironmentAccess(1.0)

	popSize := 200
	demeCount := 10
	numGens := 500

	members := make([][]population.Individual, demeCount)
	for j := 0; j < demeCount; j++ {
		members[j] = make([]population.Individual, popSize/demeCount)
		for i := 0; i < popSize/demeCount; i++ {
			members[j][i] = gp.GenerateGP(gpOpt)
		}
	}
	s := selection.ProbabilisticSelection{
		4,
		2,
	}

	pair := pairing.RandomPairing{}

	demes := make([]population.Population, demeCount)
	for i := 0; i < demeCount; i++ {
		demes[i] = population.Population{
			Members:      members[i],
			Size:         popSize / demeCount,
			Selection:    s,
			Pairing:      pair,
			FitnessTests: 4,
			TestInputs:   in,
			TestExpected: out,
			Elites:       2,
			Fitnesses:    make([]int, popSize/demeCount),
			GoalFitness:  2,
		}
	}
	dg := population.DemeGroup{
		Demes:           demes,
		MigrationChance: 0.05,
	}

	for i := 0; i < numGens; i++ {
		fmt.Println("Gen", i+1)
		stopEarly := dg.NextGeneration()
		if i == numGens-1 || stopEarly {
			for _, p := range dg.Demes {
				w, _ := p.Weights(1.0)
				fmt.Println(w)
				fmt.Println(p.Fitnesses)
				m, f := p.BestMember()
				m.Print()
				fmt.Println("Fitness: ", f)
			}
			_, f := dg.BestMember()
			fmt.Println("Best Fitness: ", f)
			fmt.Println("Generations taken: ", i+1)
			break
		}
	}

}
