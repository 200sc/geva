package lgp

import (
	"fmt"
	"goevo/env"
	"goevo/population"
	"math/rand"
)

// A linear genetic program
type LGP struct {
	Instructions []Instruction
	Mem          *env.I
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
)

const (
	LAST_WRITTEN = -1
)

var (
	gpOptions        LGPOptions
	crossover        LGPCrossover
	environment      *env.I
	memory           *env.I
	actions          []Action
	actionWeights    []float64
	cumActionWeights []float64
	fitness          FitnessFunc
)

func Init(genOpt LGPOptions, e, m *env.I, cross LGPCrossover,
	act []Action, baseActionWeight float64, f FitnessFunc) {

	actions = act

	actionWeights = make([]float64, len(actions))
	for i := range actions {
		actionWeights[i] = baseActionWeight
	}
	ResetCumActionWeights()

	environment = e
	memory = m
	fitness = f
	gpOptions = genOpt
	crossover = cross
}

func GenerateLGP(genOpt LGPOptions) *LGP {

	gp := new(LGP)

	gp.Env = environment.Copy()
	gp.Mem = memory.Copy()

	l := rand.Intn(genOpt.MaxStartActions-genOpt.MinStartActions) + genOpt.MinStartActions

	gp.Instructions = make([]Instruction, l)
	for i := range gp.Instructions {
		gp.Instructions[i] = gp.GetInstruction()
	}

	return gp
}

func (gp *LGP) Run() {
	quit_early := 300
	i := 0
	gp.lastRegister = 0
	gp.pc = 0
	// This is just to figure out who is crashing
	id := rand.Intn(1000000)
	defer func(i int) {
		if r := recover(); r != nil {
			fmt.Println("Recovered panicking id", i)
			gp.Print()
			panic("Resuming Panic")
		}
	}(id)
	nextPC := gp.pc
	for i < quit_early && nextPC < len(gp.Instructions) {
		//fmt.Println(id, nextPC, len(gp.Instructions))
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
		fmt.Println("---", i.String())
	}
	fmt.Println("MEM:")
	for i, m := range *gp.Mem {
		fmt.Println("---", i, ":", *m)
	}
	fmt.Println("LR", gp.lastRegister)
	fmt.Println("PC", gp.pc)
	fmt.Println("")
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
		gp.ValueMutate()
	case v5 < gpOptions.MemMutationChance:
		gp.MemMutate()
	}
}

func (gp *LGP) Copy() *LGP {
	gp2 := new(LGP)
	gp2.Instructions = gp.Instructions
	gp2.Mem = gp.Mem.Copy()
	gp2.Env = gp.Env.Copy()
	return gp2
}
