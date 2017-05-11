package eda

var BenchTest = And(MaxIterations(2000),
	And(TrackBest(true),
		And(Samples(100),
			And(LearningSamples(10),
				BaseValue(0.5)))))
