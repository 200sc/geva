package selection

import (
	"goevo/neural"
	"goevo/population"
	"math"
)

type StochasticUniversalSelection struct {
	ParentProportion int
	RankBased        bool
	Power            float64
}

func (sus StochasticUniversalSelection) GetParentProportion() int {
	return sus.ParentProportion
}

func (sus StochasticUniversalSelection) Select(p_p *population.Population) []neural.Network {
	p := *p_p

	fitnessChannels := p_p.Fitness()
	fitnesses := make([]int, p.Size)
	weights := make([]float64, p.Size)
	cumulativeWeights := make([]float64, p.Size)

	maxFitness := 0

	for i := 0; i < p.Size; i++ {
		v := <-fitnessChannels[i]
		if v > maxFitness {
			maxFitness = v
		}
		fitnesses[i] = v
	}

	// Transform values which are low to equivalent high
	// values on the same scale, applying the power
	// as a further bias scaling towards the best
	// individuals.
	for i := 0; i < p.Size; i++ {
		weights[i] = math.Pow(float64((fitnesses[i]*-1)+maxFitness+1), sus.Power)
	}

	cumulativeWeights[0] = weights[0]

	for i := 0; i < p.Size-1; i++ {
		cumulativeWeights[i+1] = cumulativeWeights[i] + weights[i+1]
	}

	inc := cumulativeWeights[len(cumulativeWeights)-1] / float64(p.Size/sus.ParentProportion)

	outNet := make([]neural.Network, p.Size)

	i := 0
	j := 0
	next := 0.0
	for i < len(cumulativeWeights) {
		if cumulativeWeights[i] < next {
			i++
		} else {
			outNet[j] = p.Members[i]
			j++
			next += inc
		}
	}

	return outNet
}
