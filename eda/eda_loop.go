package eda

import "bitbucket.org/StephenPatrick/goevo/env"

type EDA func(...Option) Model

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
	model := eda(opts...)
	for model.ShouldContinue() {
		model = model.Adjust(samples)
	}
}
