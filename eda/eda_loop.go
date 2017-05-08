package eda

import "bitbucket.org/StephenPatrick/goevo/env"

type EDA func(...Option) (Model, error)

type Model interface {
	BaseModel() *Base
	ShouldContinue() bool
	Adjust(samples int) Model
	// Should be generalized
	ToEnv() *env.F
}

type Fitness func(m Model) int

type Problem struct {
	length  int
	fitness Fitness
}

func Loop(eda EDA, samples int, opts ...Option) {
	model, err := eda(opts...)
	if err != nil {
		panic(err)
	}
	for model.ShouldContinue() {
		model = model.Adjust(samples)
	}
}
