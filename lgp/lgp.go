package lgp

import (
	"fmt"
	"github.com/200sc/geva/env"
	"github.com/200sc/geva/pop"
	"math/rand"
)

// A linear genetic program
type LGP struct {
	Instructions []Instruction
	Mem          *env.I
	MemStart     *env.I
	Env          *env.I
	// The last register is the last register
	// that the LGP has written to. When using
	// the special value LAST_WRITTEN, this register's
	// value is accessed by the LGP.
	lastRegister int
	pc           int
}

const (
	SPECIAL_REGISTERS = 1
	LAST_WRITTEN      = -1
)

var (
	gpOptions        Options
	crossover        LGPCrossover
	environment      *env.I
	memory           *env.I
	actions          []Action
	actionWeights    []float64
	cumActionWeights []float64
	fitness          FitnessFunc
	quit_early       int
)

func GeneratePopulation(opt interface{}, popSize int) []pop.Individual {
	gpOpt := opt.(Options)
	members := make([]pop.Individual, popSize)
	for j := 0; j < popSize; j++ {
		members[j] = GenerateLGP(gpOpt)
	}
	return members
}

func GenerateLGP(genOpt Options) *LGP {

	gp := new(LGP)

	gp.Env = environment.Copy()
	gp.Mem = memory.Copy()
	gp.MemStart = memory.Copy()

	l := rand.Intn(genOpt.MaxStartActions-genOpt.MinStartActions) + genOpt.MinStartActions

	gp.Instructions = make([]Instruction, l)
	for i := range gp.Instructions {
		gp.Instructions[i] = gp.GetInstruction()
	}

	return gp
}

func PrintActions() {
	fmt.Println(actions)
	fmt.Println(actionWeights)
}

func (gp *LGP) Run() {
	i := 0

	gp.lastRegister = 0
	gp.pc = 0
	gp.Mem = gp.MemStart.Copy()

	nextPC := gp.pc
	for i < quit_early && nextPC < len(gp.Instructions) {
		inst := gp.Instructions[nextPC]
		gp.pc++
		inst.Act.Op(gp, inst.Args...)
		i++
		nextPC = gp.pc
	}
}

func (gp *LGP) Print() {
	// Todo
	fmt.Println("Instructions:")
	for _, i := range gp.Instructions {
		fmt.Println("───", i.String())
	}
	fmt.Println("MEM:")
	for i, m := range *gp.Mem {
		fmt.Println("───", i, ":", *m)
	}
	fmt.Println("ENV:")
	for i, e := range *gp.Env {
		fmt.Println("───", i, ":", *e)
	}
	fmt.Println("LR", gp.lastRegister)
	fmt.Println("PC", gp.pc)
	fmt.Println("")
}

func (gp *LGP) CanCrossover(other pop.Individual) bool {
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
func (gp *LGP) Crossover(other pop.Individual) pop.Individual {
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
		gp.ValueMutate()
	case v5 < gpOptions.MemMutationChance:
		gp.MemMutate()
	}
}

func (gp *LGP) Copy() *LGP {
	gp2 := new(LGP)
	gp2.Instructions = gp.Instructions
	gp2.Mem = gp.Mem.Copy()
	gp2.MemStart = gp.MemStart.Copy()
	gp2.Env = gp.Env.Copy()
	return gp2
}
