package gg

import (
	"math/rand"
	"sort"

	"github.com/200sc/geva/gg/dev"
	"github.com/200sc/geva/gg/player"
	"github.com/200sc/geva/pairing"
	"github.com/200sc/geva/pop"
	"github.com/200sc/geva/selection"
)

type Instance struct {
	DevCreator    dev.Creator
	PlayerCreator player.Creator

	DevCt          int
	PlayerCt       int
	DevIterations  int
	PlayIterations int
	PlayTime       int

	Assignment func(playerCt, devCt int) [][]int
	// goal fitness?

	pop               pop.Population
	mechanics         []*dev.Mechanic
	fitnesses         []float64
	players           []player.Player
	playerAssignments [][]int
}

func (in *Instance) Loop() {
	// Create devs
	in.pop.Size = in.DevCt
	in.pop.Members = make([]pop.Individual, in.DevCt)
	in.pop.Fitnesses = make([]int, in.DevCt)
	in.pop.FitnessTests = 1
	// todo: settings for these
	in.pop.Elites = in.DevCt / 20
	in.pop.Selection = selection.DeterministicTournament{2, 2}
	in.pop.Pairing = pairing.Random{}

	in.mechanics = make([]*dev.Mechanic, in.DevCt)
	for i := 0; i < in.DevCt; i++ {
		in.pop.Members[i] = in.DevCreator.NewDev()
	}

	// Loop for dev iterations
	for i := 0; i < in.DevIterations; i++ {
		// Create players
		in.players = make([]player.Player, in.PlayerCt)
		for j := 0; j < in.PlayerCt; j++ {
			in.players[j] = in.PlayerCreator.NewPlayer()
		}

		// Assign players to devs
		assignment := in.Assignment(in.PlayerCt, in.DevCt)

		// Loop for play iterations
		for j := 0; j < in.PlayIterations; j++ {

			nextAssignment := make([][]int, in.DevCt)

			// Create Mechanics for each dev
			// This could be in previous loop as well
			for k := 0; k < in.DevCt; k++ {
				in.mechanics[k] = in.pop.Members[k].(dev.Dev).Mechanic()
				nextAssignment[k] = []int{}
			}
			// Have each player play their assigned mechanic up until PlayTime
			for k := 0; k < len(assignment); k++ {
				mch := in.mechanics[k]
				for l := 0; l < len(assignment[k]); l++ {
					p := assignment[k][l]
					enjy := in.players[p].Play(mch, in.PlayTime)

					toMove := k
					// Move players if they didn't enjoy this mechanic
					if rand.Float64() > enjy {
						toMove = rand.Intn(in.DevCt)
					}
					nextAssignment[toMove] = append(nextAssignment[toMove], p)
				}
			}
			assignment = nextAssignment
		}
		// Evaluate fitness of devs by how many players they have
		// Right now, linear-- dev with most players has fitness 1,
		// second most fitness 2, etc
		pcs := make([]PlayerCount, in.DevCt)
		for j, v := range assignment {
			pcs[j] = PlayerCount{j, len(v)}
		}
		pcss := PlayerCounts(pcs)
		sort.Sort(pcss)
		lastV := pcs[0].playerCount
		nextFitness := 1
		for j, pc := range pcss {
			if pc.playerCount != lastV {
				nextFitness++
			}
			in.pop.Members[j].(dev.Dev).SetFitness(nextFitness)
			in.pop.Fitnesses[j] = nextFitness
		}

		// Evolve dev population
		in.pop.NextGeneration()
	}
	// Evaluate results
}
