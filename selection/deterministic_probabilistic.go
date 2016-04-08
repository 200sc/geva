package selection

import (
	"goevo/neural"
	population "goevo/population"
)

type DeterministicProbabilisticSelection struct {
	ParentProportion int
}

func (dps DeterministicProbabilisticSelection) GetParentProportion() int {
	return dps.ParentProportion
}

func (dps_p *DeterministicProbabilisticSelection) Select(p_p *population.Population) []neural.Network {
	//p := *p_p

	return []neural.Network{}
}
