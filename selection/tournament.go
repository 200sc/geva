package selection

import (
	"goevo/population"
	"math/rand"
	"sort"
)

type TournamentSelection struct {
	TournamentSize int
	// 2 for 1/2, 3 for 1/3, etc
	// The remaining fraction will
	// be taken from crossover
	ParentProportion   int
	ChanceToSelectBest float64
}

func (ts TournamentSelection) GetParentProportion() int {
	return ts.ParentProportion
}

func (ts TournamentSelection) Select(p_p *population.Population) []population.Individual {
	p := *p_p

	members := make([]population.Individual, p.Size)
	// We have an arbitrary buffer here.
	// It should just effect how many goroutines can
	// simultaneously end (or all end prior to a
	// single loop pulling out values).
	selectionCh := make(chan int, 20)

	// Send off goroutines to process tournament battles
	for i := 0; i < p.Size/ts.ParentProportion; i++ {

		// Get a random set of indexes
		fighters := Sample(ts.TournamentSize, p.Size)
		fitMap := make(map[int]int)

		// Process fitness channels and map
		// fitnesses to indexes.
		for _, j := range fighters {
			fitMap[p.Fitnesses[j]] = j
		}

		// The goroutine which will pick a winner
		// of the tournament fight based on our
		// selection method's chance to pick
		// the best fitness.
		go func(fitMap map[int]int, selectionCh chan int, p float64) {

			keys := KeySet_IntInt(fitMap)
			sort.Ints(keys)

			// We take a random float and continually
			// lower the what we compare it against.
			r := rand.Float64()
			chance := 1 - p
			for i := 0; i < len(keys)-1; i++ {
				// Once our chance is less than the
				// given random number, we return
				// our current index of the sorted weights.
				if chance < r {
					selectionCh <- fitMap[keys[i]]
					return
				}
				chance -= chance * p
			}
			// On ejection we default to the last
			// index, so in some sense the last
			// index has a very minimal bias
			selectionCh <- fitMap[keys[len(keys)-1]]

		}(fitMap, selectionCh, ts.ChanceToSelectBest)
	}

	// Pull the above indexes as they are calculated
	for i := 0; i < p.Size/ts.ParentProportion; i++ {
		members[i] = p.Members[<-selectionCh]
	}
	close(selectionCh)

	return members
}
