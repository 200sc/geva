package gg

import "github.com/200sc/go-dist/floatrange"

type GGCrossover interface {
	Crossover(a, b *GG) *GG
}

var (
	acRange = floatrange.NewLinear(0, 3)
)

// Crossover Concept 1
// I don't have a better name for this concept yet
//
// Take parent A's environment
// and initial state.
//
// Pick a total number of actions
// somewhere between the number
// for parent A and parent B, and take
// half of that from parent A.
//
// Recreate the other half from the
// settings used for some actions for
// parent B.
//
// Create a new goal state from these actions.

func CrossoverOne(a, b *GG) *GG {
	c := new(GG)

	c.Environment = a.Environment
	c.Init = a.Init

	actionCount := ((len(a.Actions) + len(b.Actions)) / 2) + acRange.Poll()
	passiveCount := ((len(a.Passives) + len(b.Passives)) / 2) + acRange.Poll()

	// Hypothetically, the passive and action lists
	// for a and b are already shuffled. We shuffle the
	// end list before returning c.
	aActions := a.Actions[0 : actionCount/2]
	aPassives := a.Passives[0 : passiveCount/2]

	// Recreate bActions and bPassives . . .

	return c
}
