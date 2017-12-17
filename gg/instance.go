package gg

import (
	"fmt"
	"image/png"
	"math/rand"
	"os"
	"sort"

	"github.com/200sc/geva/gg/dev"
	"github.com/200sc/geva/gg/player"
	"github.com/200sc/geva/pairing"
	"github.com/200sc/geva/pop"
	"github.com/200sc/geva/selection"
	"github.com/200sc/geva/unique"
	"github.com/oakmound/oak"
	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/scene"
)

type Instance struct {
	DevCreator    dev.Creator
	PlayerCreator player.Creator

	DevCt           int
	PlayerCt        int
	DevIterations   int
	PlayIterations  int
	PlayTime        int
	MechanicsPerGen int

	Assignment func(playerCt, devCt int) [][]int
	// goal fitness?

	pop               pop.Population
	mechanics         []*dev.Mechanic
	fitnesses         []float64
	players           []player.Player
	playerAssignments [][]int

	Render bool
}

var (
	renderGraph = &unique.Render{}
)

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

	if in.Render {
		oak.Add(
			"uniqueness",
			func(string, interface{}) {},
			func() bool { return true },
			func() (string, *scene.Result) {
				return "uniqueness", nil
			},
		)
		go oak.Init("uniqueness")
	}

	in.mechanics = make([]*dev.Mechanic, in.DevCt)
	for i := 0; i < in.DevCt; i++ {
		in.pop.Members[i] = in.DevCreator.NewDev()
	}

	graph := unique.NewGraph(unique.MinDistance(20))

	// Loop for dev iterations
	for i := 0; i < in.DevIterations; i++ {
		fmt.Println("Iteration", i)
		// Create players
		in.players = make([]player.Player, in.PlayerCt)
		for j := 0; j < in.PlayerCt; j++ {
			in.players[j] = in.PlayerCreator.NewPlayer()
		}

		// Assign players to devs
		assignment := in.Assignment(in.PlayerCt, in.DevCt)

		totalEnjoyment := 0.0
		totalEnjoymentIt := 0.0

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
					totalEnjoyment += enjy
					totalEnjoymentIt++

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
		fmt.Println("Total Enjoyment", totalEnjoyment, "Percentage", totalEnjoyment/totalEnjoymentIt)
		fmt.Println("Evaluating fitness")
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
			lastV = pc.playerCount
			dv := in.pop.Members[j].(dev.Dev)
			thisFitness := nextFitness
			if !graph.CanAdd(dv.Mechanic()) {
				thisFitness += 10
				thisFitness *= 2
			}
			dv.SetFitness(thisFitness)
			in.pop.Fitnesses[j] = thisFitness
		}

		//for i := 0; i < in.MechanicsPerGen; i++ {
		best, _ := in.pop.BestMember()
		ok := graph.Add(dev.NewRenderMechanic(best.(dev.Dev).Mechanic()))
		if !ok {
			fmt.Println("Failed to add mechanic to uniqueness graph")
		}
		//}

		if in.Render {
			renderGraph.UnDraw()
			renderGraph.SetGraph(graph)
			render.Draw(renderGraph, 1)
		}
		// Evolve dev population
		in.pop.NextGeneration()
	}
	// Evaluate results
	fmt.Println("Best developer:")
	best, _ := in.pop.BestMember()
	fmt.Println(best.(*dev.Base))
	fmt.Println("Best mechanics:")
	for i, m := range dev.MechanicNames {
		fmt.Println(i, m.String())
	}
	if in.Render {
		rgba := oak.ScreenShot()
		f, err := os.Create("screen.png")
		if err != nil {
			dlog.Error(err)
			return
		}
		png.Encode(f, rgba)
	}
}
