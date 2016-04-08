package selection

import (
	"goevo/neural"
	population "goevo/population"
)

type ProbabilisticSelection struct {
	ParentProportion int
}

func (ps_p *ProbabilisticSelection) Select(p_p *population.Population) []neural.Network {
	//p := *p_p

	return []neural.Network{}
}
