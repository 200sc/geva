package gg

import (
	"math/rand"
	"testing"
	"time"

	"github.com/200sc/geva/cross"
	"github.com/200sc/geva/mut"
	"github.com/200sc/geva/mut/irange"
	"github.com/200sc/geva/mut/mutenv"

	"github.com/200sc/geva/mut/frange"

	"github.com/200sc/go-dist/floatrange"
	"github.com/200sc/go-dist/intrange"

	"github.com/200sc/geva/gg/dev"
	"github.com/200sc/geva/gg/player"

	. "github.com/200sc/geva/gg"
)

func TestInstanceOne(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())

	ins := &Instance{}
	// test 1
	// ins.DevCt = 80
	// ins.PlayerCt = 300
	// ins.DevIterations = 20
	// ins.PlayIterations = 8
	// ins.PlayTime = 40

	// test 2
	// ins.DevCt = 40
	// ins.PlayerCt = 600
	// ins.DevIterations = 40
	// ins.PlayIterations = 5
	// ins.PlayTime = 20

	// test 3
	ins.DevCt = 150
	ins.PlayerCt = 700
	ins.DevIterations = 5
	ins.PlayIterations = 10
	ins.PlayTime = 20
	ins.Render = true

	ins.Assignment = func(playerCt, devCt int) [][]int {
		players := make([]int, playerCt)
		for i := 0; i < playerCt; i++ {
			players[i] = i
		}
		for i := 0; i < playerCt; i++ {
			j := rand.Intn(playerCt)
			players[i], players[j] = players[j], players[i]
		}
		if playerCt < devCt {
			panic("Player count must exceed or equal dev count")
		}
		out := make([][]int, devCt)
		inc := playerCt / devCt
		p := 0
		for i := 0; i < len(out); i++ {
			out[i] = []int{}
			for j := 0; j < inc; j++ {
				out[i] = append(out[i], players[p])
				p++
			}
		}
		for p < playerCt {
			i := rand.Intn(len(out))
			out[i] = append(out[i], players[p])
			p++
		}
		return out
	}

	imut := irange.OrAny(
		.25,
		irange.Add(1),
		irange.Add(-1),
		irange.Scale(.5),
		irange.Scale(2),
	)

	fmut := frange.Or(
		frange.Add(1),
		frange.Or(
			frange.Add(-1),
			frange.Or(
				frange.Scale(.5),
				frange.Scale(2),
				.5,
			),
			.5,
		),
		.5,
	)

	imut0 := irange.And(
		imut,
		irange.EnforceMin(0),
	)

	imut1 := irange.And(
		imut,
		irange.EnforceMin(1),
	)

	imutlimit10 := irange.And(
		imut1,
		irange.EnforceMax(10),
	)

	imutlimit20 := irange.And(
		imut1,
		irange.EnforceMax(20),
	)

	fmut1 := frange.And(
		fmut,
		frange.EnforceMin(1),
	)

	ins.DevCreator = &dev.LinearCreator{
		InitSizeBottom:       floatrange.Constant(0),
		InitSizeTop:          floatrange.Constant(6),
		GoalDistanceBottom:   floatrange.Constant(2),
		GoalDistanceTop:      floatrange.Constant(10),
		GoalSizeBottom:       floatrange.Constant(1),
		GoalSizeTop:          floatrange.Constant(3),
		EnvSizeBottom:        floatrange.Constant(6),
		EnvSizeTop:           floatrange.Constant(10),
		EnvValBottom:         floatrange.Constant(-5),
		EnvValTop:            floatrange.Constant(-5),
		ActionCountBottom:    floatrange.Constant(1),
		ActionCountTop:       floatrange.Constant(4),
		ActionTypeCount:      intrange.NewLinear(3, 20),
		ActionTypeChoices:    dev.BaseActionTypes,
		ActionTypeWeights:    floatrange.NewLinear(1, 10),
		ActionStrengthBottom: floatrange.Constant(-5),
		ActionStrengthTop:    floatrange.Constant(5),
		CrossoverOptions: []dev.Crossover{
			&dev.LinearDevCrossover{
				ActionTypeCrossover: &dev.ActionModCrossover{
					TypeLengthMod: intrange.NewLinear(-1, 1),
					TypeWeightMod: floatrange.NewLinear(.5, 1.5),
					TypeWeightDef: floatrange.NewLinear(1, 4),
					StrengthMod: frange.Or(
						frange.Add(1),
						frange.Add(-1),
						.5,
					),
				},
			},
		},
		MutationOptions: []dev.DevMutation{
			dev.BasicDevMutation{
				InitSize:     imut0,
				GoalSize:     imut1,
				GoalDistance: imut1,
				EnvSize:      imutlimit20,
				EnvVal:       fmut1,
				ActionCount:  imutlimit10,
				ActionTypeWeights: mut.Or(
					mut.Or(mut.Add(.1), mut.Add(-.1), .5),
					mut.DropOut(0.5), .99),
				ActionStrengths: fmut,
			},
		},
	}

	ins.PlayerCreator = &player.IntEnvCreator{
		ExpectedTime:    floatrange.NewLinear(0.5, .95),
		ExpectedFitness: intrange.NewLinear(1, 10),
		PopSizeBottom:   intrange.NewLinear(3, 10),
		PopSizeTop:      intrange.NewLinear(12, 20),
		ActionCtBottom:  intrange.NewLinear(1, 3),
		ActionCtTop:     intrange.NewLinear(4, 20),
		Mutators: []mutenv.I{
			mutenv.OnAll(
				mut.Or(
					mut.Add(1),
					mut.Add(-1),
					.5,
				),
			),
		},
		Crosses: []cross.I{
			cross.IPointCrossover{1},
			cross.IPointCrossover{2},
		},
	}

	ins.Loop()
}
