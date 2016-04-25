package selection

import (
	//"fmt"
	"goevo/neural"
	"goevo/population"
	"sort"
)

type GreedySelection struct {
	ParentProportion int
}

func (gs GreedySelection) GetParentProportion() int {
	return gs.ParentProportion
}

// I'll be honest, I don't know what 'greed' is supposed
// to mean in this exact circumstance of evolutionary selection.
// I saw it referenced once.
// I'm assuming here that it means picking the top x members of the
// population where x is the proportion of members who are parents
// in the new generation.
func (gs GreedySelection) Select(p_p *population.Population) []neural.Network {
	p := *p_p

	fitMap := make(map[int][]int)
	members := make([]neural.Network, p.Size)

	// Send off goroutines to process tournament battles
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
	for i < p.Size/gs.ParentProportion {
		for k := 0; k < len(fitMap[keys[j]]); k++ {
			if i >= p.Size/gs.ParentProportion {
				return members
			}
			members[i] = p.Members[fitMap[keys[j]][k]]
			i++
		}
		j++
	}
	return members
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
