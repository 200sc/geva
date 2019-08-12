package gp

import (
	"math/rand"
	"strconv"

	"github.com/200sc/geva/env"
)

// The effect of a given action is internally defined
// by the action, and GPs need to learn what actions
// are useful, when.
// The first slice here represents how many arguments
// each action takes.
type Actions [][]Action
type Action struct {
	Op   Operator
	Name string
}

// This is severely limited so as to
// weigh all actions the same. We want to be able
// to give certain actions different weights for
// being picked!
var (
	OneArgActions = []Action{
		{getEnv, "env"},
		//{neg, "neg"},
		//{pow2, "pow2"},
		//{pow3, "pow3"},
	}
	TwoArgActions = []Action{
		{add, "add"},
		//{subtract, "sub"},
		{multiply, "mult"},
		{divide, "div"},
		//{pow, "pow"},
		//{mod, "mod"},
		//{ifRand, "rand?"},
		//{do2, "do2"},
		// Tree based GPs NEED something like this
		// to compete with LGPs, but this method does
		// not achieve the desired results.
		//{whilePositive, "while+"},
	}
	ThreeArgActions = []Action{
		{neZero, "!0?"},
		//{isPositive, "+?"},
		//{ifRand, "rand?"},
		//{do3, "do3"},
	}
	ZeroArgActions = []Action{
		//{randv, "rand"},
		{zero, "0"},
		{one, "1"},
		{two, "2"},
		//{three, "3"},
		//{four, "4"},
		//{five, "5"},
		//{six, "6"},
		//{seven, "7"},
		//{eight, "8"},
		//{nine, "9"},
	}
	Pow8Actions = [][]Action{
		{
			{randv, "rand"},
			{zero, "0"},
			{one, "1"},
			{two, "2"},
		},
		{
			{neg, "neg"},
			{getEnv, "env"},
		},
		{
			{add, "add"},
			{subtract, "sub"},
			{multiply, "multiply"},
		},
		{
			{isPositive, "+?"},
			{ifRand, "rand?"},
			{neZero, "!0?"},
		},
	}
	BaseActions = [][]Action{
		ZeroArgActions,
		OneArgActions,
		TwoArgActions,
		ThreeArgActions,
	}
	TartarusActions = [][]Action{
		{
			{randv, "rand"},
			{zero, "0"},
			{one, "1"},
			{two, "2"}},
		{
			{neg, "neg"}},
		{
			{add, "add"},
			{subtract, "sub"},
			{multiply, "multiply"},
			{do2, "do2"}},
		{
			{isPositive, "+?"},
			{ifRand, "rand?"},
			{do3, "do3"},
			{neZero, "!0?"}},
	}
)

func getAction(args ...int) (action Action, children int) {
	// Make a choice out of the options
	// and treat that choice as an index
	// as if the available arrays were attached end-on-end

	choice := weightedChooseOne(CalculateCumulativeActionWeights(args...))

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
	choice := weightedChooseOne(cumZeroActionWeights)
	action = actions[0][choice]
	return
}

func getNonZeroAction() (action Action, children int) {
	// Make a choice out of the options
	// and treat that choice as an index
	// as if the available arrays were attached end-on-end
	choice := weightedChooseOne(cumActionWeights)

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

func getEnv(gp *GP, xs ...*Node) int {
	index := Eval(xs[0])
	if index >= len(*gp.Env) {
		index = len(*gp.Env) - 1
	}
	if index < 0 {
		index = 0
	}
	return *(*gp.Env)[index]
}

func setEnv(gp *GP, xs ...*Node) int {
	index := Eval(xs[1])
	if index >= len(*gp.Env) {
		index = len(*gp.Env) - 1
	}
	if index < 0 {
		index = 0
	}
	v := Eval(xs[0])
	*(*gp.Env)[index] = v
	return v
}

func AddStorage(spaces int, baseWeight float64) {
	storageActions := make([]Action, spaces)
	storageWeights := make([]float64, spaces)

	for i := 0; i < spaces; i++ {
		j := i
		storageActions[i] = Action{
			func(gp *GP, xs ...*Node) int {
				v := Eval(xs[0])
				*(*gp.Mem)[j] = v
				return v
			},
			"sav" + strconv.Itoa(j)}
		storageWeights[i] = baseWeight
	}

	actions[1] = append(actions[1], storageActions...)
	actionWeights[1] = append(actionWeights[1], storageWeights...)
	cumActionWeights = CalculateCumulativeActionWeights(1, 2, 3)

	memActions := make([]Action, spaces)
	memWeights := make([]float64, spaces)

	for i := 0; i < spaces; i++ {
		j := i
		memActions[i] = Action{
			func(gp *GP, nothing ...*Node) int {
				return *(*gp.Mem)[j]
			},
			"mem" + strconv.Itoa(j)}
		memWeights[i] = baseWeight
	}

	actions[0] = append(actions[0], memActions...)
	actionWeights[0] = append(actionWeights[0], memWeights...)
	cumZeroActionWeights = CalculateCumulativeActionWeights(0)

	if memory == nil {
		memory = new(env.I)
	}
	*memory = make([]*int, spaces)
}

func CalculateCumulativeActionWeights(args ...int) []float64 {
	weightCount := 0
	for _, arg := range args {
		weightCount += len(actions[arg])
	}
	w := make([]float64, weightCount)
	w[len(w)-1] = actionWeights[1][0]
	k := len(w) - 2
	for i, arg := range args {
		for j := 0; j < len(actions[arg]); j++ {
			if i+j != 0 {
				w[k] = w[k+1] + actionWeights[arg][j]
				k--
			}
		}
	}
	return w
}

// weightedChooseOne returns a single index from the weights given
// at a rate relative to the magnitude of each weight. It expects
// the input to be in the form of RemainingWeights, cumulative with
// the total at index 0.
func weightedChooseOne(remainingWeights []float64) int {
	totalWeight := remainingWeights[0]
	choice := rand.Float64() * totalWeight
	i := len(remainingWeights) / 2
	start := 0
	end := len(remainingWeights) - 1
	for {
		if remainingWeights[i] < choice {
			if remainingWeights[i-1] < choice {
				end = i
				i = (start + end) / 2
			} else {
				return i - 1
			}
		} else {
			if i != len(remainingWeights)-1 && remainingWeights[i+1] > choice {
				start = i

				i = (start + end) / 2
				if (start+end)%2 == 1 {
					i++
				}
			} else {
				return i
			}
		}
	}
}

// Returns whether the modification was successful (whether
// an action with that name was found)
func ModifyActionWeight(action string, newWeight float64) bool {
	for i, tier := range actions {
		for j, a := range tier {
			if a.Name == action {
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
