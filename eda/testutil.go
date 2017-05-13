package eda

import (
	"math/rand"
	"time"

	"bitbucket.org/StephenPatrick/goevo/mut"
	"github.com/200sc/go-dist/floatrange"
)

// BenchTest is a set of Options each benchmark test should go through to
// guarantee that different methods are compared to eachother fairly. It
// may get more complex as time goes on.
var BenchTest = And(
	func(m Model) { Seed() },
	MaxIterations(2000),
	TrackBest(true),
	Samples(100),
	LearningSamples(10),
	BaseValue(0.5),
	TrackFitnessRuns(true),
	MutationRate(.15),
	FMutator(
		mut.And(
			mut.Or(
				mut.Or(mut.Add(.1), mut.Add(-.1), .5),
				mut.DropOut(0.5), .99),
			floatrange.NewLinear(0.0, 1.0).EnforceRange)))

func Seed() {
	rand.Seed(time.Now().UTC().UnixNano())
}
