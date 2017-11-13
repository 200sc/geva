package lgp

import (
	"github.com/200sc/geva/alg"
	"math/rand"
	"strconv"
)

type Instruction struct {
	Act  Action
	Args []int
}

func (i *Instruction) String() string {
	s := i.Act.Name
	s += ":"
	for _, a := range i.Args {
		s += strconv.Itoa(a) + " "
	}
	return s
}

// The effect of a given action is internally defined
// by the action, and GPs need to learn what actions
// are useful, when.
// Unlike Tree-GPs, LGPs do not care about the number of
// arguments their actions take unless mutating or generating them.
type Action struct {
	Op   Operator
	Name string
	Args int
}

var (
	BaseActions = []Action{
		{neg, "neg", 1},
		//{pow2, "pow2", 2},
		//{pow3, "pow3", 2},
		{add, "add", 3},
		{subtract, "sub", 3},
		{multiply, "mult", 3},
		//{divide, "div", 3},
		//{pow, "pow", 3},
		//{mod, "mod", 3},
		{bnez, "bnez", 1},
		{bgz, "bgz", 1},
		{jmp, "jmp", 1},
		{randv, "rand", 1},
		{zero, "0", 1},
		{one, "1", 1},
		//{two, "2", 1},
		//{three, "3", 1},
		//{four, "4", 1},
		//{five, "5",1 },
		//{six, "6", 1},
		//{seven, "7", 1},
		//{eight, "8", 1},
		//{nine, "9", 1},
	}
	MinActions = []Action{
		{bnez, "bnez", 1},
		{bgz, "bgz", 1},
		{jmp, "jmp", 1},
	}
	EnvActions = []Action{
		{getEnv, "env", 2},
		//{setEnv, "setEnv", 2},
		{envLen, "envLen", 1},
	}
	PowSumActions = []Action{
		{getEnv, "env", 2},
		{add, "add", 3},
		{multiply, "mult", 3},
		{divide, "div", 3},
		{zero, "0", 1},
		{one, "1", 1},
		{two, "2", 1},
	}
	TartarusActions = []Action{
		{neg, "neg", 1},
		{add, "add", 3},
		{subtract, "sub", 3},
		{multiply, "mult", 3},
		{bnez, "bnez", 1},
		{bgz, "bgz", 1},
		{jmp, "jmp", 1},
		{randv, "rand", 1},
		{zero, "0", 1},
		{one, "1", 1},
		{two, "2", 1},
	}
	SortActions = []Action{
		{bnez, "bnez", 1},
		{bgz, "bgz", 1},
		{bgeq, "bgeq", 2},
		{jmp, "jmp", 1},
		{getEnv, "env", 2},
		{envLen, "envLen", 1},
		{divide, "div", 3},
		{subtract, "sub", 3},
		{zero, "0", 1},
		{one, "1", 1},
		{memSwap, "swap", 2},
	}
)

func (gp *LGP) GetInstruction() Instruction {
	act := getAction()
	return Instruction{
		act,
		getArgs(act.Args, len(*gp.Mem)+SPECIAL_REGISTERS),
	}
}

func getAction() Action {
	return actions[alg.CumWeightedChooseOne(cumActionWeights)]
}

func getArgs(argCount int, limit int) []int {
	args := make([]int, argCount)
	for i := range args {
		// Todo: distribute these numbers non-linearly
		args[i] = rand.Intn(limit) - SPECIAL_REGISTERS
	}
	return args
}

func ResetCumActionWeights() {
	cumActionWeights = make([]float64, len(actions))
	cumActionWeights[0] = actionWeights[0]
	for i := 1; i < len(actionWeights); i++ {
		cumActionWeights[i] = cumActionWeights[i-1] + actionWeights[i]
	}
}
