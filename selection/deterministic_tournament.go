package selection

import (
	"goevo"
	"goevo/neural"
	"math"
)

type DeterministicTournamentSelection struct {
	tournamentSize     int
	selectedProportion int
}

func (dts_p *DeterministicTournamentSelection) Select(p_p *goevo.Population) []neural.Network {
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
	for i := 0; i < p.Size/ts.selectedProportion; i++ {

		// Get a random set of indexes
		fighters := Sample(ts.tournamentSize, p.Size)
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
		members[i] = p.Members[fitMap[bestFitness]]
	}

	return members
}
