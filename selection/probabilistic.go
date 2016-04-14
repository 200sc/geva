package selection

import (
	"goevo/neural"
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

func (ps ProbabilisticSelection) Select(p_p *population.Population) []neural.Network {
	p := *p_p

	weights, _ := p.Weights(ps.Power)

	maxWeight := 0.0
	for _, w := range weights {
		if w < maxWeight {
			maxWeight = w
		}
	}

	next := 0
	outNet := make([]neural.Network, p.Size)

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
