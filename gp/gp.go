package gp

import (
	"goevo/env"
	"goevo/population"
	"math/rand"
)

// The principal Individual implementation for the gp package
type GP struct {
	First *Node
	Env   *env.I
	Mem   *env.I
	Nodes int
}

var (
	gpOptions            GPOptions
	crossover            GPCrossover
	environment          *env.I
	memory               *env.I
	actions              Actions
	actionWeights        [][]float64
	cumActionWeights     []float64
	cumZeroActionWeights []float64
	fitness              FitnessFunc
)

func Init(genOpt GPOptions, e *env.I, cross GPCrossover,
	act [][]Action, baseActionWeight float64, f FitnessFunc) {

	environment = e
	actions = act
	actionWeights = make([][]float64, len(actions))
	for i, tier := range actions {
		actionWeights[i] = make([]float64, len(tier))
		for j := range tier {
			actionWeights[i][j] = baseActionWeight
		}
	}
	cumZeroActionWeights = CalculateCumulativeActionWeights(0)
	cumActionWeights = CalculateCumulativeActionWeights(1, 2, 3)
	fitness = f
	gpOptions = genOpt
	crossover = cross
}

func GenerateGP(genOpt GPOptions) *GP {

	// Eventually we'll do something
	// with creation types here

	gp := new(GP)
	gp.Env = environment.Copy()
	gp.Mem = memory.Copy()
	a, children := getNonZeroAction()
	gp.First = &Node{
		make([]*Node, children),
		a,
		gp,
	}
	gp.Nodes = gp.First.GenerateTree(genOpt.MaxStartDepth, genOpt.MaxNodeCount)
	return gp
}

func (gp *GP) Print() {
	// We aren't going to print the environment,
	// because GPs are usually going to be printed
	// in large sets, and because environments are
	// usually going to be constant between GPs.
	gp.First.Print("", true)
}

func (gp *GP) CanCrossover(other population.Individual) bool {
	switch other.(type) {
	default:
		return false
	case *GP:
		return true
	}
}

func (gp *GP) Crossover(other population.Individual) population.Individual {
	return crossover.Crossover(gp, other.(*GP))
}

func (gp *GP) Fitness(input, expected [][]float64) int {
	return fitness(gp, input, expected)
}

func (gp *GP) Mutate() {
	if rand.Float64() < gpOptions.SwapMutationChance {
		gp.SwapMutate()
	}
	if rand.Float64() < gpOptions.ShrinkMutationChance {
		gp.ShrinkMutate()
	}
	gp.Nodes = gp.First.Size()
}
