package gp

import (
	"math/rand"
)

// The effect of a given action is internally defined
// by the action, and GPs need to learn what actions
// are useful, when.
// The first slice here represents how many arguments
// each action takes.
type Actions [][]func(*GP, ...Node)

var (
	OneArgActions = []func(*GP, ...Node) int{
		neg,
	}
	TwoArgActions = []func(*GP, ...Node) int{
		add,
		subtract,
		multiply,
		divide,
		pow,
		mod,
		ifRand,
	}
	ThreeArgActions = []func(*GP, ...Node) int{
		neZero,
		isPositive,
		ifRand,
	}
	ZeroArgActions = []func(*GP, ...Node) int{
		rand,
		one,
		two,
		three,
		four,
		five,
		six,
		seven,
		eight,
		nine,
	}
	BaseActions = [][]func(*GP, ...Node) int{
		ZeroArgActions,
		OneArgActions,
		TwoArgActions,
		ThreeArgActions,
	}
)

func getAction(args ...int) (action func(*GP, ...Node)) {

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
	choice := rand.IntN(choices)

	for i := 0; i < len(args); i++ {
		if len(actions[args[i]]) < choice {
			action = actions[args[i]][choice]
			break
		}
		choice -= len(action[args[i]])
	}
	return
}

func getNonZeroAction() (action func(*GP, ...Node)) {

	var choices int
	for i := 1; i < len(actions); i++ {
		choices += len(actions[i])
	}
	// Make a choice out of the options
	// and treat that choice as an index
	// as if the available arrays were attached end-on-end
	choice := rand.IntN(choices)

	for i := 1; i < len(actions); i++ {
		if len(actions[i]) < choice {
			action = actions[args[i]][choice]
			break
		}
		choice -= len(action[args[i]])
	}
	return
}
