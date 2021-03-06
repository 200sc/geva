package selection

import (
	"github.com/200sc/geva/pop"
	"math"
)

type DeterministicTournament struct {
	TournamentSize   int
	ParentProportion int
}

func (dts DeterministicTournament) GetParentProportion() int {
	return dts.ParentProportion
}

func (dts DeterministicTournament) Select(p *pop.Population) []pop.Individual {
	// We move as much initialization down here as we can,
	// because we expect the above goroutines to be the
	// most expensive time sink in this function.
	members := make([]pop.Individual, p.Size)

	// Send off goroutines to process tournament battles
	for i := 0; i < p.Size/dts.ParentProportion; i++ {

		// Get a random set of indexes
		fighters := Sample(dts.TournamentSize, p.Size)
		fitMap := make(map[int]int)

		// Process fitness channels and map
		// fitnesses to indexes.
		bestFitness := math.MaxInt32
		for _, j := range fighters {
			fitMap[p.Fitnesses[j]] = j
			if p.Fitnesses[j] < bestFitness {
				bestFitness = p.Fitnesses[j]
			}
		}
		// Here's the point of difference between deterministic and non-deterministic
		// tournament selection-- deterministic will always pick the most fit fighter
		// as the winner of each round. non-deterministic has the built-in variable
		// where it might not pick the most fit fighter.
		members[i] = p.Members[fitMap[bestFitness]]
	}

	return members
}
