package eda

// BenchTest is a set of Options each benchmark test should go through to
// guarantee that different methods are compared to eachother fairly. It
// may get more complex as time goes on.
var BenchTest = And(
	MaxIterations(2000),
	TrackBest(true),
	Samples(100),
	LearningSamples(10),
	BaseValue(0.5),
	TrackFitnessRuns(true))
