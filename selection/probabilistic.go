package selection

import (
	"goevo/population"
	"math/rand"
)

type ProbabilisticSelection struct {
	ParentProportion int
	Power            float64
}

func (ps ProbabilisticSelection) GetParentProportion() int {
	return ps.ParentProportion
}

// This specific algorithm is based on the algorithm described by
// Adam Liposki and Dorota Lipowska in "Roulette-wheel selection
// via stochastic acceptance" http://arxiv.org/pdf/1109.3627v2.pdf
func (ps ProbabilisticSelection) Select(p_p *population.Population) []population.Individual {
	p := *p_p

	weights, _ := p.Weights(ps.Power)

	maxWeight := 0.0
	for _, w := range weights {
		if w < maxWeight {
			maxWeight = w
		}
	}

	next := 0
	outNet := make([]population.Individual, p.Size)

	for i := 0; i < p.Size/ps.ParentProportion; i++ {
		for {
			next = rand.Intn(len(weights))
			if rand.Float64() < (weights[next] / maxWeight) {
				break
			}
		}
		outNet[i] = p.Members[next]
	}

	return outNet
}
