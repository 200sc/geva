package dev

import (
	"image/color"
	"math"
	"strconv"

	"github.com/200sc/geva/unique"
	"github.com/200sc/go-dist/colorrange"
	"github.com/oakmound/oak/render"
)

var (
	_ unique.Node       = &Mechanic{}
	_ render.Renderable = &RenderMechanic{}
	_ unique.RenderNode = &RenderMechanic{}
)

var (
	MechanicNames  = map[int]*Mechanic{}
	nextName       int
	MechanicLookup = map[*Mechanic]*render.Composite{}
	mechColorRange = colorrange.NewLinear(color.RGBA{0, 0, 0, 255}, color.RGBA{255, 255, 255, 255})
)

func NextMechanicName() int {
	// not concurrency safe
	nextName++
	return nextName
}

type RenderMechanic struct {
	*render.Composite
	*Mechanic
}

func NewRenderMechanic(mch *Mechanic) *RenderMechanic {
	cmp, ok := MechanicLookup[mch]
	if !ok {
		sp := render.NewColorBox(10, 10, mechColorRange.Poll())

		nm := "m" + strconv.Itoa(NextMechanicName())
		str := render.DefFont().NewStrText(nm, 1, 1).ToSprite()

		cmp = render.NewComposite([]render.Modifiable{sp, str})
		MechanicLookup[mch] = cmp
	}
	return &RenderMechanic{
		cmp,
		mch,
	}
}

func (rmch *RenderMechanic) GetR() render.Renderable {
	return rmch.Composite
}

func (mch *Mechanic) MechanicDistance(mch2 *Mechanic) float64 {
	dist := 0.0
	mch.Reset()
	mch2.Reset()
	for i := 0; i < mch.Environment.Len(); i++ {
		if i >= mch2.Environment.Len() {
			dist++
		} else {
			dist += math.Abs(mch.Environment.Get(i) - mch2.Environment.Get(i))
		}
	}
	// Both the goal distance from a -> b and b -> a are half counted
	for i, v := range mch.Goal {
		if v2, ok := mch2.Goal[i]; ok {
			dist += math.Abs(v-v2) / 2
		} else {
			dist += .5
		}
	}
	for i, v := range mch2.Goal {
		if v2, ok := mch.Goal[i]; ok {
			dist += math.Abs(v-v2) / 2
		} else {
			dist += .5
		}
	}
	// actions can't be compared
	dist += math.Abs(float64(len(mch.Actions) - len(mch2.Actions)))
	return dist
}
