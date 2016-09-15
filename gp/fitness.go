package gp

import (
	"math"
)

type FitnessFunc func(gp *GP, inputs, outputs [][]float64) int

// An example fitness function which treats
// the output as a environment to compare
// a modified environment by the GP to.
func EnvFitness(g *GP, inputs, outputs [][]float64) int {
	fitness := 1
	for i, envDiff := range inputs {
		g.env = environment.New(envDiff)
		Eval(g.first)
		fitness += g.env.Diff(outputs[i])
	}
	return fitness
}

// An example fitness which treats the
// output of the GP as a value to compare
// against the single expected output
func OutputFitness(g *GP, inputs, outputs [][]float64) int {
	fitness := 1
	for i, envDiff := range inputs {
		g.env = environment.New(envDiff)
		out := Eval(g.first)
		fitness += int(math.Abs(float64(out - int(outputs[i][0]))))
	}
	return fitness
}
