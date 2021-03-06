package unique

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/200sc/geva/cross"
	"github.com/200sc/geva/env"
	"github.com/200sc/geva/mut"
	"github.com/200sc/geva/pop"
	"github.com/200sc/go-dist/floatrange"
)

// EnvInd is a wrapper around a float env that can be a member
// of a population
type EnvInd struct {
	*env.F
	mutChance        float64
	mutChanceMutator mut.FloatMutator
	mutator          mut.FloatMutator
	crossover        cross.F
	*Graph
}

var (
	zeroToOne = floatrange.NewLinear(0, 1).EnforceRange
)

// NewEnvInd initializes a EnvInd
func NewEnvInd(size int, baseVal float64, g *Graph) *EnvInd {
	return &EnvInd{
		F:         env.NewF(size, baseVal),
		mutChance: .1,
		mutChanceMutator: mut.And(
			mut.Or(
				mut.DropOut(.1),
				mut.Add(.01),
				.05),
			zeroToOne,
		),
		mutator: mut.Or(
			mut.Or(mut.Add(.4), mut.Add(-.4), .5),
			mut.Or(mut.Add(3), mut.Add(-3), .5),
			.5),
		crossover: cross.FPoint{NumPoints: 2},
		Graph:     g,
	}
}

func (ei *EnvInd) Fitness(input, expected [][]float64) int {
	// calculate all pairs distance
	// fitness is how different that is from the graph's all pairs distance
	pts := make([][2]float64, ei.F.Len()/2)
	for i := 0; i < ei.F.Len(); i += 2 {
		pts[i/2] = [2]float64{ei.Get(i), ei.Get(i + 1)}
	}
	distError := 0.0
	for i := 0; i < len(pts); i++ {
		pt1 := pts[i]
		n1 := ei.Graph.nodes[i]
		for j := i + 1; j < len(pts); j++ {
			n2 := ei.Graph.nodes[j]

			realDist, ok := n1.Distance(n2)
			if ok {
				pt2 := pts[j]
				dist := distance(pt1[0], pt2[0], pt1[1], pt2[1])

				distError += math.Abs(realDist - dist)
			}
		}
	}
	return int(distError) + 1
}

func distance(x1, x2, y1, y2 float64) float64 {
	return math.Sqrt(
		math.Pow(x1-x2, 2) +
			math.Pow(y1-y2, 2))
}

func (ei *EnvInd) Mutate() {
	ei.F = ei.F.Mutate(ei.mutChance, ei.mutator)
	ei.mutChance = ei.mutChanceMutator(ei.mutChance)
	rnd := rand.Float64()
	if rnd < 0.025 {
		ei.crossover = cross.FPoint{NumPoints: 2}
	} else if rnd < 0.5 {
		ei.crossover = cross.FAverageCrossover{AWeight: .5}
	}
	// could mutate env mutator

}

func (ei *EnvInd) Crossover(other pop.Individual) pop.Individual {
	if ei2, ok := other.(*EnvInd); ok {
		f := ei.F.Copy()
		mutator := ei.mutator
		if rand.Float64() < .5 {
			mutator = ei2.mutator
		}
		crossover := ei.crossover
		if rand.Float64() < .5 {
			crossover = ei2.crossover
		}
		mutChanceMutator := ei.mutChanceMutator
		if rand.Float64() < .5 {
			mutChanceMutator = ei2.mutChanceMutator
		}
		return &EnvInd{
			F:                f,
			mutChance:        ei.mutChance + ei2.mutChance/2,
			mutChanceMutator: mutChanceMutator,
			mutator:          mutator,
			crossover:        crossover,
			Graph:            ei.Graph,
		}
	}
	return ei
}

// CanCrossover always returns false for a EnvInd
func (ei *EnvInd) CanCrossover(other pop.Individual) bool {
	_, ok := other.(*EnvInd)
	return ok
}

// Print prints a EnvInd
func (ei *EnvInd) Print() {
	fmt.Println(ei.F)
}
