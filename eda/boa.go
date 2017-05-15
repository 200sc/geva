package eda

import (
	"bitbucket.org/StephenPatrick/goevo/env"
	"bitbucket.org/StephenPatrick/goevo/pop"
)

type BOA struct {
	Base
	P *pop.Population
}

func (boa *BOA) Adjust() Model {
	selected := boa.SelectLearning(boa.P)
	// Construct a network based on selected
	bn := NewBayesNet(MemberEnvs(selected))
	// Sample children from the network
	samples := bn.Sample(int(boa.learningRate * float64(boa.samples)))
	// Replace worst fitnesses in boa.P with children
	boa.ReplaceLowFitnesses(boa.P, samples)
	return boa
}

func BOAModel(opts ...Option) (Model, error) {
	var err error
	boa := new(BOA)
	boa.Base, err = DefaultBase(opts...)
	boa.P = boa.Pop()
	return boa, err
}

type BayesNet struct{}

func NewBayesNet(samples []*env.F) *BayesNet {
	return nil
}

func (bn *BayesNet) Sample(n int) []*env.F {
	samples := make([]*env.F, n)
}
