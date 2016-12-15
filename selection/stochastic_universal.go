package selection

import (
	"goevo/pop"
)

type StochasticUniversal struct {
	ParentProportion int
	// Rank based doesn't do anything
	// What it would do would it would
	// take a network's weight as its
	// ranking in a sort
	RankBased bool
	Power     float64
}

func (sus StochasticUniversal) GetParentProportion() int {
	return sus.ParentProportion
}

func (sus StochasticUniversal) Select(p_p *pop.Population) []pop.Individual {
	p := *p_p

	_, cumulativeWeights := p.Weights(sus.Power)

	inc := cumulativeWeights[len(cumulativeWeights)-1] / float64(p.Size/sus.ParentProportion)

	outNet := make([]pop.Individual, p.Size)

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
