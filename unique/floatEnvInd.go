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
	"github.com/oakmound/oak/alg/floatgeom"
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
		mutChance: .01,
		mutChanceMutator: mut.And(
			mut.Or(
				mut.DropOut(.01),
				mut.Add(.01),
				.05),
			zeroToOne,
		),
		mutator: mut.Or(
			mut.Or(mut.Add(.1), mut.Add(-.1), .5),
			mut.Or(mut.Add(1), mut.Add(-1), .5),
			.7),
		crossover: cross.FPointCrossover{NumPoints: 2},
		Graph:     g,
	}
}

func (ei *EnvInd) Fitness(input, expected [][]float64) int {
	// calculate all pairs distance
	// fitness is how different that is from the graph's all pairs distance
	dist := 0.0
	pts := make([]floatgeom.Point2, ei.F.Len()/2)
	for i := 0; i < ei.F.Len(); i += 2 {
		pts[i/2] = floatgeom.Point2{ei.Get(i), ei.Get(i + 1)}
	}
	for i := 0; i < len(pts); i++ {
		pt1 := pts[i]
		for j := i + 1; j < len(pts); j++ {
			pt2 := pts[j]
			dist += pt1.Distance(pt2)
		}
	}
	return int(math.Abs(ei.allPairsDistance-dist)) + 1
}

func (ei *EnvInd) Mutate() {
	ei.F = ei.F.Mutate(ei.mutChance, ei.mutator)
	ei.mutChance = ei.mutChanceMutator(ei.mutChance)
	rnd := rand.Float64()
	if rnd < 0.025 {
		ei.crossover = cross.FPointCrossover{NumPoints: 2}
	} else if rnd < 0.5 {
		ei.crossover = cross.FAverageCrossover{AWeight: .5}
	}
	// could mutate env mutator

}

// Crossover is  NOP on a EnvInd
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
		return &EnvInd{
			F:         f,
			mutChance: ei.mutChance + ei2.mutChance/2,
			mutator:   mutator,
			crossover: crossover,
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
