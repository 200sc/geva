package lgp

import (
	"fmt"
	"goevo/env"
	"goevo/population"
)

// A linear genetic program
type LGP struct {
	Actions      []Action
	Env          *env.I
	lastRegister int
	pc           int
}

var (
	gpOptions        LGPOptions
	crossover        LGPCrossover
	environment      *env.I
	actions          []Action
	actionWeights    [][]float64
	cumActionWeights []float64
	fitness          FitnessFunc
)

func Init(genOpt LGPOptions, e *env.I, cross LGPCrossover,
	act []Action, baseActionWeight float64, f FitnessFunc) {

	actions = act

	actionWeights = make([]float64, len(actions))
	cumActionWeights = make([]float64, len(actions))
	for i := range actions {
		actionWeights[i] = baseActionWeight
	}
	cumActionWeights[0] = actionWeights[0]
	for i := 1; i < len(actions); i++ {
		cumActionWeights[i] = actionWeights[i] + cumActionWeights[i-1]
	}

	environment = e
	fitness = f
	gpOptions = genOpt
	crossover = cross
}

func GenerateLGP(genOpt LGPOptions) *LGP {

	gp = new(LGP)

	return gp
}

func (gp *LGP) Print() {
	fmt.Println("A LGP")
	// Todo
}

func (gp *LGP) CanCrossover(other population.Individual) bool {
	switch other.(type) {
	default:
		return false
	case *LGP:
		return true
	}
}

// Crossover types
// (multi)point crossover
// uniform crossover
// these should be brought out and used for all list-like structures
func (gp *LGP) Crossover(other population.Individual) population.Individual {
	return crossover.Crossover(gp, other.(*LGP))
}

func (gp *LGP) Fitness(input, expected [][]float64) int {
	return fitness(gp, input, expected)
}

// Mutation types:
// swap mutate, swapping instructions at two locations
// value mutate, changing the values baked into an instruction
// shrink/expand mutate, removing or adding instructions from random locations
// environment mutate, changing the initial environment values
func (gp *LGP) Mutate() {
	v := rand.Float64()
	v2 := v - gpOptions.SwapMutationChance
	v3 := v2 - gpOptions.ShrinkMutationChance
	v4 := v3 - gpOptions.ExpandMutationChance
	v5 := v4 - gpOptions.ValueMutationChance
	switch {
	case v < gpOptions.SwapMutationChance:
		gp.SwapMutate()
	case v2 < gpOptions.ShrinkMutationChance:
		gp.ShrinkMutate()
	case v3 < gpOptions.ExpandMutationChance:
		gp.ExpandMutate()
	case v4 < gpOptions.ValueMutationChance:
		gp.ValueMUtate()
	case v5 < gpOptions.EnvMutationChance:
		gp.EnvMutate()
	}
}
