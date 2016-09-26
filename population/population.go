package population

import (
	"math"
	"math/rand"
	"sort"
)

type Population struct {
	Members      []Individual
	Size         int
	Selection    SelectionMethod
	Pairing      PairingMethod
	FitnessTests int
	TestInputs   [][]float64
	TestExpected [][]float64
	Elites       int
	Fitnesses    []int
	LowFitness   int
	MaxFitness   int
	GoalFitness  int
}

// This will change as more things take place
// in a generation. Selection, Crossover, and Mutation
// are granted.
func (p_p *Population) NextGeneration() bool {
	p := *p_p
	// The number of parents in the next generation
	parentSize := p.Size / p.Selection.GetParentProportion()

	p.GenerateFitness()
	if p.LowFitness <= p.GoalFitness {
		return true
	}
	elites := p.GetElites()
	nextGen := p.Selection.Select(&p)

	// Ensure that the elites (the best members)
	// stay in the next generation
	for i, elite := range elites {
		nextGen[i+parentSize] = nextGen[i]
		nextGen[i] = elite
	}
	parentSize += p.Elites

	p.Members = nextGen
	pairs := p.Pairing.Pair(&p, parentSize)

	// i does not start at 0,
	// but pairs, sensibly, does.
	pairIndex := 0

	// crossover pairs for children in the next generation.
	for i := parentSize; i < len(nextGen); i++ {
		n1 := p.Members[pairs[pairIndex][0]]
		n2 := p.Members[pairs[pairIndex][1]]
		p.Members[i] = n1.Crossover(n2)
		pairIndex++
	}

	// Mutate. The elites are not subject to mutation.
	for i := p.Elites; i < len(p.Members); i++ {
		p.Members[i].Mutate()
	}

	*p_p = p
	return false
}

func (p_p *Population) GenerateFitness() {
	channels := make([]chan int, p_p.Size)

	for i := 0; i < p_p.Size; i++ {
		channels[i] = make(chan int)

		go func(n Individual, ch chan int, inputs [][]float64, expected [][]float64, tests int) {
			if len(inputs) == tests {
				ch <- n.Fitness(inputs, expected)
				return
			}
			in := make([][]float64, tests)
			out := make([][]float64, tests)
			for i := 0; i < tests; i++ {
				j := rand.Intn(len(inputs))
				in[i] = inputs[j]
				out[i] = expected[j]
			}
			ch <- n.Fitness(in, out)
		}((p_p.Members[i]), channels[i], p_p.TestInputs, p_p.TestExpected, p_p.FitnessTests)
	}

	p_p.LowFitness = math.MaxInt32
	p_p.MaxFitness = 0

	for i := 0; i < p_p.Size; i++ {
		v := <-channels[i]
		close(channels[i])
		if v < p_p.LowFitness {
			p_p.LowFitness = v
		} else if v > p_p.MaxFitness {
			p_p.MaxFitness = v
		}
		p_p.Fitnesses[i] = v
	}
}

func (p_p *Population) GetElites() []Individual {
	p := *p_p

	fitMap := make(map[int][]int)
	elites := make([]Individual, p.Elites)

	for i := 0; i < p.Size; i++ {
		f := p.Fitnesses[i]
		if v, ok := fitMap[f]; ok {
			fitMap[f] = append(v, i)
		} else {
			fitMap[f] = []int{i}
		}
	}

	keys := KeySet_Int_SlInt(fitMap)
	sort.Ints(keys)
	i := 0
	j := 0
	for i < p.Elites {
		for k := 0; k < len(fitMap[keys[j]]); k++ {
			if i >= p.Elites {
				return elites
			}
			elites[i] = p.Members[fitMap[keys[j]][k]]
			i++
		}
		j++
	}
	return elites
}

func KeySet_Int_SlInt(m map[int][]int) []int {
	keys := make([]int, len(m))

	i := 0
	for k := range m {
		keys[i] = k
		i++
	}

	return keys
}

func (p_p *Population) Weights(power float64) ([]float64, []float64) {
	p := *p_p
	fitnesses := p.Fitnesses
	maxFitness := p.MaxFitness

	weights := make([]float64, len(fitnesses))
	cumulativeWeights := make([]float64, len(fitnesses))

	// Transform values which are low to equivalent high
	// values on the same scale, applying the power
	// as a further bias scaling towards the best
	// individuals.
	for i := 0; i < len(fitnesses); i++ {
		weights[i] = math.Pow(float64((fitnesses[i]*-1)+maxFitness+1), power)
	}

	cumulativeWeights[0] = weights[0]

	for i := 0; i < len(fitnesses)-1; i++ {
		cumulativeWeights[i+1] = cumulativeWeights[i] + weights[i+1]
	}

	return weights, cumulativeWeights
}

func (p *Population) BestMember() (Individual, int) {
	w, _ := p.Weights(1.0)
	maxWeight := 0.0
	maxIndex := 0
	for i, v := range w {
		if v > maxWeight {
			maxWeight = v
			maxIndex = i
		}
	}
	return p.Members[maxIndex], p.Fitnesses[maxIndex]
}

func (p_p *Population) Print() {
	for _, v := range p_p.Members {
		v.Print()
	}
}

// Used as Generic helpers for populations

type SelectionMethod interface {
	Select(p *Population) []Individual
	GetParentProportion() int
}

type PairingMethod interface {
	Pair(p *Population, populated int) [][]int
}
