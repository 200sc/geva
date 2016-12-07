package lgp

import (
	"goevo/algorithms"
	"goevo/env"
	"strconv"
)

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
		//{subtract, "sub", 3},
		{multiply, "mult", 3},
		//{divide, "div", 3},
		//{pow, "pow", 3},
		//{mod, "mod", 3},
		{bnez, "bnez", 1},
		{bgz, "bgz", 1},
		{jmp, "jmp", 1},
		{randv, "rand", 1},
		//{one, "1", 1},
		//{two, "2", 1},
		//{three, "3", 1},
		//{four, "4", 1},
		//{five, "5",1 },
		//{six, "6", 1},
		//{seven, "7", 1},
		//{eight, "8", 1},
		//{nine, "9", 1},
	}
)

func getAction() Action {
	return actions[algorithms.CumWeightedChooseOne(cumActionWeights)]
}
