package neural

import (
	"sort"
	"rand"
)

type Population []Network

type PopulationOptions struct {
	generationOptions *NetworkGenerationOptions
	size int
	selection *PopulationSelectionMethod
	crossover *PopulationCrossoverMethod
	testInputs [][]bool
	testExpected [][]bool
}

func (p_p *Population) NextGeneration(pOpt_p *PopulationOptions) *Population {
	p := *p_p

	nextGen := p.selection.Select(p_p, pOpt_p)
	return p.crossover.Crossover(nextGen) 
}

type PopulationSelectionMethod interface {
	Select(p_p *Population, pOpt_p *PopulationOptions) *Population
}

type PopulationCrossoverMethod interface {
	Crossover(p_p *Population) *Population
}

// Stores options for Greedy selection
type GreedySelection struct {
	selectedProportion int
}

func (gs_p *GreedySelection) Select(p_p *Population, pOpt_p *PopulationOptions) *Population {
	//gs := *gs_p
	//p := *p_p

	return p_p
}

type DeterministicTournamentSelection struct {
	tournamentSize int
	selectedProportion int
}

type TournamentSelection struct {
	tournamentSize int
	// 2 for 1/2, 3 for 1/3, etc
	// The remaining fraction will
	// be taken from crossover
	selectedProportion int
	chanceToSelectBest float64
}

func (ts_p *TournamentSelection) Select(p_p *Population, pOpt_p *PopulationOptions) *Population {
	p := *p_p
	pOpt := *pOpt_p


	// Send off goroutines to calculate the population members' fitnesses
	channels := make([]chan int, pOpt.size) 
	for i := 0; i < pOpt.size; i++ {

		go func(n *Network, ch chan int, inputs []bool, expected []bool) {
			ch <- (*n).Fitness(inputs, expected)

		}(&(p[i]), channels[i], pOpt.testInputs, pOpt.testExpected)
	}

	// We move as much initialization down here as we can,
	// because we expect the above goroutines to be the
	// most expensive time sink in this function. 
	ts := *ts_p
	fitnesses := make([]int, pOpt.size)
	newPopulation := make([]Network pOpt.size)
	// We have an arbitrary buffer here. 
	// It should just effect how many goroutines can 
	// simultaneously end (or all end prior to a
	// single loop pulling out values). 
	selectionCh := make(chan int, 20)

	// Send off goroutines to process tournament battles
	for i := 0; i < pOpt.size / ts.selectedProportion; i++ {

		// Get a random set of indexes
		fighters := Sample(ts.tournamentSize, pOpt.size)
		fitMap := make(map[int]int)

		// Process fitness channels and map
		// fitnesses to indexes. 
		for _, j := range fighters {
			// The slice will be initialized to zero
			// and we return 1 as optimal fitness,
			// so this is a check for initialization.
			if fitnesses[j] == 0 {
				fitnesses[j] <- channels[j]
				close(channels[j])
			}
			fitMap[fitnesses[j]] = j
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
			chance := 1-p			
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

		}(fitMap, selectionCh, ts.chanceToSelectBest)
	}

	// Pull the above indexes as they are calculated
	for i := 0; i < pOpt.size / ts.selectedProportion; i++ {
		newPopulation[i] = p[<-selectionCh]
	}
	close(selectionCh)

	return &newPopulation
}

type ProbabilisticSelection struct {
	selectedProportion int
}

func (ps_p *ProbabilisticSelection) Select(p_p *Population, pOpt_p *PopulationOptions) *Population {
	//ps := *ps_p
	//p := *p_p

	return p_p
}
