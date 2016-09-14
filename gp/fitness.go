package gp

import (
	"math"
)

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
	for i, env := range inputs {
		g.env = environment.New(envDiff)
		out := Eval(g.first)
		fitness += math.Abs(out - int(output[i][0]))
	}
	return fitness
}
