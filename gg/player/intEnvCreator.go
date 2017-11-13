package player

import (
	"math/rand"

	"github.com/200sc/geva/cross"
	"github.com/200sc/geva/mut/mutenv"
	"github.com/200sc/go-dist/floatrange"
	"github.com/200sc/go-dist/intrange"
)

type IntEnvCreator struct {
	ExpectedTime    floatrange.Range
	ExpectedFitness intrange.Range

	PopSizeBottom intrange.Range
	PopSizeTop    intrange.Range

	ActionCtBottom intrange.Range
	ActionCtTop    intrange.Range

	Mutators []mutenv.I
	Crosses  []cross.I
}

func (iec *IntEnvCreator) NewPlayer() Player {
	return &IntEnvPlayer{
		expectedTime:    iec.ExpectedTime.Poll(),
		expectedFitness: iec.ExpectedFitness.Poll(),
		popSize:         intrange.NewLinear(iec.PopSizeBottom.Poll(), iec.PopSizeTop.Poll()),
		actionCt:        intrange.NewLinear(iec.ActionCtBottom.Poll(), iec.ActionCtTop.Poll()),
		Mutator:         iec.Mutators[rand.Intn(len(iec.Mutators))],
		Cross:           iec.Crosses[rand.Intn(len(iec.Crosses))],
	}
}
