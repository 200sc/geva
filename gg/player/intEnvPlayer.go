package player

import (
	"fmt"
	"math/rand"

	"github.com/200sc/geva/cross"
	"github.com/200sc/geva/env"
	"github.com/200sc/geva/gg/dev"
	"github.com/200sc/geva/mut/mutenv"
	"github.com/200sc/geva/pairing"
	"github.com/200sc/geva/pop"
	"github.com/200sc/geva/selection"
	"github.com/200sc/go-dist/intrange"
)

type IntEnvPlayer struct {
	expectedTime    float64
	expectedFitness int

	popSize  intrange.Range
	actionCt intrange.Range
	pop      pop.Population

	Mutator mutenv.I
	Cross   cross.I
	mch     *dev.Mechanic
}

func (iip *IntEnvPlayer) Play(mch *dev.Mechanic, playTime int) float64 {
	iip.mch = mch

	iip.pop.Size = iip.popSize.Poll()
	iip.pop.Members = make([]pop.Individual, iip.pop.Size)
	iip.pop.Fitnesses = make([]int, iip.pop.Size)
	iip.pop.FitnessTests = 1
	// todo: settings for these
	iip.pop.Elites = 1
	iip.pop.Selection = selection.DeterministicTournament{2, 2}
	iip.pop.Pairing = pairing.Random{}
	iip.pop.GoalFitness = iip.expectedFitness

	for i := 0; i < iip.pop.Size; i++ {
		iip.pop.Members[i] = iip.RandomIntInd(len(mch.Actions), iip.actionCt.Poll())
	}

	i := 0
	for ; i < playTime; i++ {
		if iip.pop.NextGeneration() {
			break
		}
	}
	expectedGeneration := int(float64(playTime) * iip.expectedTime)
	_, bestFitness := iip.pop.BestMember()
	if bestFitness >= iip.expectedFitness {
		if i <= expectedGeneration {
			// All expectations met
			return 1.0
		}
		// todo: work on these hardcoded values
		// game beaten, but not as fast as desired
		return 0.5
	}
	// game not beaten
	return 0.1
}

func (iip *IntEnvPlayer) Fitness(e *env.I) int {
	iip.mch.Reset()
	for _, j := range *e {
		iip.mch.Actions[*j]()
		// todo: consider a forgiving impl
		// forgiving would mean that if the goal is reached, the player knows immediately
	}
	return iip.mch.MechFitness(iip.mch)
}

type IntInd struct {
	*env.I
	player *IntEnvPlayer
}

func (iip *IntEnvPlayer) RandomIntInd(actionLength int, size int) *IntInd {
	e := env.NewI(size, 0)
	for i := 0; i < size; i++ {
		*(*e)[i] = rand.Intn(actionLength)
	}
	return &IntInd{e, iip}
}

func (ii *IntInd) Fitness(input, expected [][]float64) int {
	return ii.player.Fitness(ii.I)
}

func (ii *IntInd) Mutate() {
	ii.player.Mutator(ii.I)
}

func (ii *IntInd) Crossover(other pop.Individual) pop.Individual {
	return &IntInd{
		ii.player.Cross.Crossover(ii.I, other.(*IntInd).I),
		ii.player,
	}
}

func (ii *IntInd) CanCrossover(other pop.Individual) bool {
	_, ok := other.(*IntInd)
	return ok
}

func (ii *IntInd) Print() {
	fmt.Println(ii.I)
}
