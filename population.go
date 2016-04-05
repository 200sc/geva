package neural

import (
	"sort"
	"rand"
)

type Population struct {
	members []Network
	generationOptions *NetworkGenerationOptions
	size int
	selection *SelectionMethod
	crossover *CrossoverMethod
	testInputs [][]bool
	testExpected [][]bool
}

func (p_p *Population) NextGeneration() *Population {
	p := *p_p

	nextGen := p.selection.Select(p_p)
	p.members = p.crossover.Crossover(nextGen) 
	return &p
}

type SelectionMethod interface {
	Select(p_p *Population) []Network
}

type GreedySelection struct {
	selectedProportion int
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

type ProbabilisticSelection struct {
	selectedProportion int
}

func (p_p *Population) Fitness() []chan int {
	p := *p_p

	channels := make([]chan int, p.size) 

	for i := 0; i < p.size; i++ {

		go func(n *Network, ch chan int, inputs []bool, expected []bool) {
			ch <- (*n).Fitness(inputs, expected)

		}(&(p[i]), channels[i], p.testInputs, p.testExpected)
	}

	return channels
}

func (dts_p *DeterministicTournamentSelection) Select(p_p *Population) []Network {
	p := *p_p


	// Send off goroutines to calculate the population members' fitnesses
	fitnessChannels := p_p.Fitness(p_p)

	// We move as much initialization down here as we can,
	// because we expect the above goroutines to be the
	// most expensive time sink in this function. 
	ts := *ts_p
	fitnesses := make([]int, p.size)
	members := make([]Network, p.size)

	// Send off goroutines to process tournament battles
	for i := 0; i < p.size / ts.selectedProportion; i++ {

		// Get a random set of indexes
		fighters := Sample(ts.tournamentSize, p.size)
		fitMap := make(map[int]int)

		// Process fitness channels and map
		// fitnesses to indexes.
		bestFitness := math.MaxInt32
		for _, j := range fighters {
			// The slice will be initialized to zero
			// and we return 1 as optimal fitness,
			// so this is a check for initialization.
			if fitnesses[j] == 0 {
				fitnesses[j] <- fitnessChannels[j]
				close(channels[j])
			}
			fitMap[fitnesses[j]] = j
			if fitnesses[j] < bestFitness {
				bestFitness = fitnesses[j]
			}
		}
		members[i] = fitMap[bestFitness]
	}

	return members	
}

func (ts_p *TournamentSelection) Select(p_p *Population) []Network {
	p := *p_p


	// Send off goroutines to calculate the population members' fitnesses
	fitnessChannels := p_p.Fitness(p_p)

	// We move as much initialization down here as we can,
	// because we expect the above goroutines to be the
	// most expensive time sink in this function. 
	ts := *ts_p
	fitnesses := make([]int, p.size)
	members := make([]Network, p.size)
	// We have an arbitrary buffer here. 
	// It should just effect how many goroutines can 
	// simultaneously end (or all end prior to a
	// single loop pulling out values). 
	selectionCh := make(chan int, 20)

	// Send off goroutines to process tournament battles
	for i := 0; i < p.size / ts.selectedProportion; i++ {

		// Get a random set of indexes
		fighters := Sample(ts.tournamentSize, p.size)
		fitMap := make(map[int]int)

		// Process fitness channels and map
		// fitnesses to indexes. 
		for _, j := range fighters {
			// The slice will be initialized to zero
			// and we return 1 as optimal fitness,
			// so this is a check for initialization.
			if fitnesses[j] == 0 {
				fitnesses[j] <- fitnessChannels[j]
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
	for i := 0; i < p.size / ts.selectedProportion; i++ {
		members[i] = p[<-selectionCh]
	}
	close(selectionCh)

	return members
}

func (ps_p *ProbabilisticSelection) Select(p_p *Population) []Network {
	//p := *p_p

	return p_p
}

func (gs_p *GreedySelection) Select(p_p *Population) []Network {
	//gs := *gs_p
	//p := *p_p

	return p_p
}

type CrossoverMethod interface {
	Crossover(p_p *Population) []Network
}