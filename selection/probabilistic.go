package selection

import (
	"github.com/200sc/geva/pop"
	"math/rand"
)

type Probabilistic struct {
	ParentProportion int
	Power            float64
}

func (ps Probabilistic) GetParentProportion() int {
	return ps.ParentProportion
}

// This specific algorithm is based on the algorithm described by
// Adam Liposki and Dorota Lipowska in "Roulette-wheel selection
// via stochastic acceptance" http://arxiv.org/pdf/1109.3627v2.pdf
func (ps Probabilistic) Select(p_p *pop.Population) []pop.Individual {
	p := *p_p

	weights, _ := p.Weights(ps.Power)

	maxWeight := 0.0
	for _, w := range weights {
		if w < maxWeight {
			maxWeight = w
		}
	}

	next := 0
	outNet := make([]pop.Individual, p.Size)

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
