package goevo

import (
	"testing"

	"github.com/200sc/geva/neural"
	"github.com/200sc/geva/pairing"
	"github.com/200sc/geva/pop"
	"github.com/200sc/geva/selection"
	"github.com/200sc/go-dist/intrange"
)

func TestNNRun(t *testing.T) {

	Seed()

	nngOpt := neural.NetworkGenerationOptions{
		NetworkMutationOptions: neural.NetworkMutationOptions{
			WeightOptions: neural.FloatMutationOptions{
				MutChance:     0.20,
				MutMagnitude:  0.05,
				MutRange:      20,
				ZeroOutChance: 0.01,
			},
			ColumnOptions: neural.ColumnGenerationOptions{
				MinSize:           3,
				MaxSize:           4,
				DefaultAxonWeight: 0.5,
			},
			ActivatorOptions:        neural.AllActivators,
			NeuronReplacementChance: 0.05,
			NeuronAdditionChance:    0.00,
			WeightSwapChance:        0.05,
			ColumnRemovalChance:     0.00,
			ColumnAdditionChance:    0.00,
			NeuronMutationChance:    0.10,
			ActivatorMutationChance: 0.01,
		},
		MinColumns:    1,
		MaxColumns:    2,
		MaxInputs:     3,
		MaxOutputs:    1,
		BaseMutations: 20,
	}

	in := [][]float64{
		{3.0, 2.0, 0.0},
		{10.0, 20.0, 10.0},
		{2.0, 100.0, 1.0},
		{0.0, 0.0, 50.0},
		{10.0, 1.0, 1.0},
	}
	out := make([][]float64, len(in))
	for i, f := range in {
		out[i] = []float64{(f[0] + f[1] + f[2]) * 3}
	}

	neural.Init(
		nngOpt,
		neural.AverageCrossover{2},
		neural.AbsFitness,
	)

	RunSuite(
		[]TestCase{{in, out, "x3Test"}},
		4,
		200,
		100000,
		nngOpt,
		neural.GeneratePopulation,
		[]pop.SMethod{selection.Probabilistic{3, 1.7}},
		[]pop.PMethod{pairing.Random{}},
		2.0,
		intrange.NewLinear(2, 3),
		0.1,
		"ENN",
	)
}
