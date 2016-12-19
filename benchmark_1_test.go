package goevo

// TGPs need the ability to use an arbitrary number count as nodes
// to be able to easily access all elements of this environment
// func TestSuiteOne(t *testing.T) {

// 	Seed()

// 	suites := make([]TestSuite, 0)

// 	gpOpt := gp.Options{
// 		MaxNodeCount:         250,
// 		MaxStartDepth:        6,
// 		MaxDepth:             12,
// 		SwapMutationChance:   0.10,
// 		ShrinkMutationChance: 0.05,
// 	}

// 	suites = append(suites, TestSuite{
// 		[]TestCase{MultiplyMatrixTestCase()},
// 		25,
// 		3000,
// 		100000,
// 		gpOpt,
// 		gp.GeneratePopulation,
// 		gp.InitOptions{
// 			gpOpt,
// 			env.NewI(1, 0),
// 			gp.PointCrossover{},
// 			gp.BaseActions,
// 			1.0,
// 			gp.MatchMemFitness,
// 			20,
// 			1.0,
// 		},
// 		gp.OptInit,
// 		[]pop.SMethod{selection.DeterministicTournament{2, 3}},
// 		[]pop.PMethod{pairing.Random{}},
// 		1,
// 		alg.LinearIntRange{4, 6},
// 		0.05,
// 		"TGP",
// 	})

// 	lgpOpt := lgp.Options{
// 		MinActionCount:  10,
// 		MaxActionCount:  200,
// 		MaxStartActions: 80,
// 		MinStartActions: 20,

// 		SwapMutationChance:   0.15,
// 		ValueMutationChance:  0.15,
// 		ShrinkMutationChance: 0.10,
// 		ExpandMutationChance: 0.10,
// 		MemMutationChance:    0.00,
// 	}

// 	l_actions := lgp.BaseActions
// 	l_actions = append(l_actions, lgp.EnvActions...)

// 	suites = append(suites, TestSuite{
// 		[]TestCase{MultiplyMatrixTestCase()},
// 		25,
// 		3000,
// 		100000,
// 		lgpOpt,
// 		lgp.GeneratePopulation,
// 		lgp.InitOptions{
// 			lgpOpt,
// 			env.NewI(5, 0),
// 			env.NewI(20, 0),
// 			lgp.PointCrossover{3},
// 			l_actions,
// 			1.0,
// 			lgp.MatchMemFitness,
// 			600,
// 		},
// 		lgp.OptInit,
// 		[]pop.SMethod{
// 			selection.Probabilistic{3, 2},
// 			selection.Probabilistic{2, 2},
// 			selection.DeterministicTournament{2, 3},
// 			selection.DeterministicTournament{3, 3},
// 			selection.Tournament{4, 3, 0.5},
// 		},
// 		[]pop.PMethod{pairing.Random{}},
// 		1,
// 		alg.LinearIntRange{1, 10},
// 		0.10,
// 		"LGP",
// 	})

// 	nngOpt := neural.NetworkGenerationOptions{
// 		NetworkMutationOptions: neural.NetworkMutationOptions{
// 			WeightOptions: neural.FloatMutationOptions{
// 				MutChance:     0.20,
// 				MutMagnitude:  2.0,
// 				MutRange:      60,
// 				ZeroOutChance: 0.01,
// 			},
// 			ColumnOptions: neural.ColumnGenerationOptions{
// 				MinSize:           3,
// 				MaxSize:           4,
// 				DefaultAxonWeight: 0.5,
// 			},
// 			ActivatorOptions:        neural.AllActivators,
// 			NeuronReplacementChance: 0.05,
// 			NeuronAdditionChance:    0.00,
// 			WeightSwapChance:        0.05,
// 			ColumnRemovalChance:     0.00,
// 			ColumnAdditionChance:    0.00,
// 			NeuronMutationChance:    0.10,
// 			ActivatorMutationChance: 0.01,
// 		},
// 		MinColumns:    10,
// 		MaxColumns:    11,
// 		MaxInputs:     18,
// 		MaxOutputs:    9,
// 		BaseMutations: 20,
// 	}

// 	suites = append(suites, TestSuite{
// 		[]TestCase{MultiplyMatrixTestCase()},
// 		4,
// 		200,
// 		100000,
// 		nngOpt,
// 		neural.GeneratePopulation,
// 		neural.InitOptions{
// 			nngOpt,
// 			neural.AverageCrossover{2},
// 			neural.MatchFitness(0.25),
// 		},
// 		neural.OptInit,
// 		[]pop.SMethod{selection.DeterministicTournament{3, 3}},
// 		[]pop.PMethod{pairing.Random{}},
// 		2.0,
// 		alg.LinearIntRange{1, 4},
// 		0.1,
// 		"ENN",
// 	})

// 	RunTestSuites(suites)
// }
