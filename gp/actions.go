package gp

import (
	"math/rand"
	"strconv"
)

// The effect of a given action is internally defined
// by the action, and GPs need to learn what actions
// are useful, when.
// The first slice here represents how many arguments
// each action takes.
type Actions [][]Action
type Operator func(*GP, ...*Node) int

// I am dumb and this is severly limited so as to
// weigh all actions the same. We want to be able
// to give certain actions different weights for
// being picked!
var (
	OneArgActions = []Action{
		{neg, "neg"},
	}
	TwoArgActions = []Action{
		{add, "add"},
		{subtract, "sub"},
		{multiply, "mult"},
		{divide, "div"},
		{pow, "pow"},
		{mod, "mod"},
		{ifRand, "rand?"},
	}
	ThreeArgActions = []Action{
		{neZero, "!0?"},
		{isPositive, "+?"},
		{ifRand, "rand?"},
	}
	ZeroArgActions = []Action{
		{randv, "rand"},
		{one, "1"},
		{two, "2"},
		{three, "3"},
		{four, "4"},
		{five, "5"},
		{six, "6"},
		{seven, "7"},
		{eight, "8"},
		{nine, "9"},
	}
	BaseActions = [][]Action{
		ZeroArgActions,
		OneArgActions,
		TwoArgActions,
		ThreeArgActions,
	}
)

func getAction(args ...int) (action Action, children int) {

	var choices int
	// For each given arg, representing an amount of args to
	// an action, increase our total number of choices by the
	// number of actions available at that amount.
	for i := 0; i < len(args); i++ {
		choices += len(actions[args[i]])
	}
	// Make a choice out of the options
	// and treat that choice as an index
	// as if the available arrays were attached end-on-end
	choice := rand.Intn(choices)

	for i := 0; i < len(args); i++ {
		if len(actions[args[i]]) > choice {
			action = actions[args[i]][choice]
			children = i
			break
		}
		choice -= len(actions[args[i]])
	}
	return
}

func getNonZeroAction() (action Action, children int) {

	var choices int
	for i := 1; i < len(actions); i++ {
		choices += len(actions[i])
	}
	// Make a choice out of the options
	// and treat that choice as an index
	// as if the available arrays were attached end-on-end
	choice := rand.Intn(choices)

	for i := 1; i < len(actions); i++ {
		if len(actions[i]) > choice {
			action = actions[i][choice]
			children = i
			break
		}
		choice -= len(actions[i])
	}
	return
}

func AddEnvironmentAccess() {
	envActions := make([]Action, len(environment))
	for i := range environment {
		envActions[i] = Action{
			func(gp *GP, nothing ...*Node) int {
				return *(*gp.env)[i]
			},
			"env" + strconv.Itoa(i)}
	}
	actions[0] = append(actions[0], envActions...)
}

// If we like the above pattern, we can
// also add a similar AddEnvironmentChanging
// function.
