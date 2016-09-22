package gp

import (
	//"fmt"
	"goevo/algorithms"
	"strconv"
)

// The effect of a given action is internally defined
// by the action, and GPs need to learn what actions
// are useful, when.
// The first slice here represents how many arguments
// each action takes.
type Actions [][]Action
type Action struct {
	op   Operator
	name string
}

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
	// Make a choice out of the options
	// and treat that choice as an index
	// as if the available arrays were attached end-on-end

	choice := algorithms.CumWeightedChooseOne(CalculateCumulativeActionWeights(args...))

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

func getZeroAction() (action Action) {
	// Make a choice out of the options
	// and treat that choice as an index
	// as if the available arrays were attached end-on-end
	choice := algorithms.CumWeightedChooseOne(cumZeroActionWeights)
	action = actions[0][choice]
	return
}

func getNonZeroAction() (action Action, children int) {
	// Make a choice out of the options
	// and treat that choice as an index
	// as if the available arrays were attached end-on-end
	choice := algorithms.CumWeightedChooseOne(cumActionWeights)

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

func AddEnvironmentAccess(baseWeight float64) {
	envActions := make([]Action, len(environment))
	envWeights := make([]float64, len(environment))
	for i := range environment {
		envActions[i] = Action{
			func(gp *GP, nothing ...*Node) int {
				return *(*gp.env)[i]
			},
			"env" + strconv.Itoa(i)}
		envWeights[i] = baseWeight
	}
	actions[0] = append(actions[0], envActions...)
	actionWeights[0] = append(actionWeights[0], envWeights...)
	cumZeroActionWeights = CalculateCumulativeActionWeights(0)
}

// If we like the above pattern, we can
// also add a similar AddEnvironmentChanging
// function.

func CalculateCumulativeActionWeights(args ...int) []float64 {
	weightCount := 0
	for i := 0; i < len(args); i++ {
		weightCount += len(actions[args[i]])
	}
	w := make([]float64, weightCount)
	w[0] = actionWeights[1][0]
	k := 1
	for i := 0; i < len(args); i++ {
		for j := 0; j < len(actions[args[i]]); j++ {
			if i+j != 0 {
				w[k] = w[k-1] + actionWeights[args[i]][j]
				k++
			}
		}
	}
	return w
}

// Returns whether the modification was successful (whether
// an action with that name was found)
func ModifyActionWeight(action string, newWeight float64) bool {
	for i, tier := range actions {
		for j, a := range tier {
			if a.name == action {
				actionWeights[i][j] = newWeight
				if i != 0 {
					cumActionWeights = CalculateCumulativeActionWeights(1, 2, 3)
				}
				return true
			}
		}
	}
	return false
}
