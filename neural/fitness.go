package neural

import "math"

type FitnessFunc func(nn *Network, inputs, outputs [][]float64) int

func AbsFitness(nn *Network, inputs, expected [][]float64) int {
	fitness := 1.0
	for i := range inputs {
		output := nn.Run(inputs[i])
		for j := range expected[i] {
			fitness += math.Abs(output[j] - expected[i][j])
		}
	}
	return int(math.Ceil(fitness))
}

func MatchFitness(tolerance float64) FitnessFunc {
	return func(nn *Network, inputs, expected [][]float64) int {
		fitness := 1
		for i := range inputs {
			output := nn.Run(inputs[i])
			for j := range expected[i] {
				if math.IsNaN(output[j]) {
					fitness++
				}
				diff := math.Abs(output[j] - expected[i][j])
				if diff > tolerance {
					fitness++
				}
			}
		}
		return fitness
	}
}
