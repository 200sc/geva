package dev

import (
	"github.com/200sc/geva/mut"
	"github.com/200sc/geva/mut/frange"
	"github.com/200sc/geva/mut/irange"
)

type DevMutation interface {
	Mutate(*Base)
}

type BasicDevMutation struct {
	InitSize          irange.Mutator
	GoalSize          irange.Mutator
	GoalDistance      irange.Mutator
	EnvSize           irange.Mutator
	EnvVal            frange.Mutator
	ActionCount       irange.Mutator
	ActionTypeWeights mut.FloatMutator
	ActionStrengths   frange.Mutator
}

func (bdv BasicDevMutation) Mutate(d *Base) {
	d.InitSize = bdv.InitSize(d.InitSize)
	d.GoalSize = bdv.GoalSize(d.GoalSize)
	d.GoalDistance = bdv.GoalDistance(d.GoalDistance)
	d.EnvSize = bdv.EnvSize(d.EnvSize)
	d.EnvVal = bdv.EnvVal(d.EnvVal)
	d.ActionCount = bdv.ActionCount(d.ActionCount)
	for i, v := range d.ActionTypeWeights {
		d.ActionTypeWeights[i] = bdv.ActionTypeWeights(v)
	}
	for i, v := range d.ActionStrengths {
		d.ActionStrengths[i] = bdv.ActionStrengths(v)
	}
}
