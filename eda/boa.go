package eda

import (
	"math"

	"bitbucket.org/StephenPatrick/goevo/eda/stat"
	"bitbucket.org/StephenPatrick/goevo/env"
	"bitbucket.org/StephenPatrick/goevo/pop"
)

// BOA represents the Bayesian Optimization Algorithm
type BOA struct {
	Base
	P *pop.Population
}

// Mutate on a BOA mutate's the BOA's population
func (boa *BOA) Mutate() {
	for _, m := range boa.P.Members {
		m.(*EnvInd).F.Mutate(boa.mutationRate, boa.fmutator)
	}
}

// Adjust on a BOA produces a bayes net from selections of the
// BOA's population, then replaces part of its population with
// samples from the bayes net.
func (boa *BOA) Adjust() Model {
	selected := boa.SelectLearning(boa.P)
	// Update boa.F to represent selected (Perhaps not needed if we changed
	// the fitness method)
	boa.F.SetAll(0.0)
	for _, s := range selected {
		boa.F.AddF(s.(*EnvInd).F)
	}
	boa.F.Divide(float64(len(selected)))
	// Construct a network based on selected
	envs := MemberEnvs(selected)
	bn := NewBayesNet(envs)
	// Sample children from the network
	samples := bn.Sample(envs, int(boa.learningRate*float64(boa.samples)))
	// Replace worst fitnesses in boa.P with children
	boa.ReplaceLowFitnesses(boa.P, samples)
	return boa
}

// BOAModel returns a BOA EDA
func BOAModel(opts ...Option) (Model, error) {
	var err error
	boa := new(BOA)
	boa.Base, err = DefaultBase(opts...)
	boa.P = boa.Pop()
	return boa, err
}

// This is some tricky garbage the original code I based this off of
// does that makes the code difficult to understand.

func countEdges(samples []*env.F, indices []int) []int {
	counts := make([]int, int(math.Pow(2, float64(len(indices)))))
	revIndices := make([]int, len(indices))
	for i, v := range indices {
		revIndices[len(indices)-(1+i)] = v
	}
	for _, s := range samples {
		j := 0.0
		for i, v := range revIndices {
			if s.Get(v) == 1.0 {
				j += math.Pow(2, float64(i))
			}
		}
		counts[int(j)]++
	}
	return counts
}

// Because this is based on the above tricky garbage, it shared
// some of the unfortunate difficulty

func k2(i int, candidates []int, samples []*env.F) float64 {
	edgeCounts := countEdges(samples, append(candidates, i))
	total := 1.0
	for j := 0; j < len(edgeCounts)/2; j++ {
		a1, a2 := edgeCounts[j*2], edgeCounts[(j*2)+1]
		total *= (1.0 / (float64(stat.Factorial(a1+a2) + 1))) *
			float64(stat.Factorial(a1)) * float64(stat.Factorial(a2))
	}
	return total
}
