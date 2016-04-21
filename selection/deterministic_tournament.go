package selection

import (
	"goevo/neural"
	"goevo/population"
	"math"
)

type DeterministicTournamentSelection struct {
	TournamentSize   int
	ParentProportion int
}

func (dts_p *DeterministicTournamentSelection) Select(p_p *population.Population) []neural.Network {
	p := *p_p

	// Send off goroutines to calculate the population members' fitnesses
	fitnessChannels := p_p.Fitness()

	// We move as much initialization down here as we can,
	// because we expect the above goroutines to be the
	// most expensive time sink in this function.
	ts := *dts_p
	fitnesses := make([]int, p.Size)
	members := make([]neural.Network, p.Size)

	// Send off goroutines to process tournament battles
	for i := 0; i < p.Size/ts.ParentProportion; i++ {

		// Get a random set of indexes
		fighters := Sample(ts.TournamentSize, p.Size)
		fitMap := make(map[int]int)

		// Process fitness channels and map
		// fitnesses to indexes.
		bestFitness := math.MaxInt32
		for _, j := range fighters {
			// The slice will be initialized to zero
			// and we return 1 as optimal fitness,
			// so this is a check for initialization.
			if fitnesses[j] == 0 {
				fitnesses[j] = <-fitnessChannels[j]
				close(fitnessChannels[j])
			}
			fitMap[fitnesses[j]] = j
			if fitnesses[j] < bestFitness {
				bestFitness = fitnesses[j]
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
