package unique

import (
	"fmt"
	"math/rand"

	"github.com/oakmound/oak/alg/floatgeom"

	"github.com/200sc/geva/pairing"
	"github.com/200sc/geva/pop"
	"github.com/200sc/geva/selection"

	"github.com/oakmound/oak/render"
)

var (
	_ render.Renderable = &Render{}
)

type Render struct {
	*Graph
	*render.CompositeR
}

func (r *Render) UnDraw() {
	if r.CompositeR == nil {
		return
	}
	r.CompositeR.UnDraw()
}

func (r *Render) SetGraph(g *Graph) {
	r.Graph = g
	// there's necessarily going to be error if the nodes
	// represent more than two variables. We move elements to
	// minimize error, up to some number of tries.
	//
	// Is this a genetic algorithm problem
	// yes
	p := pop.Population{}
	p.Members = make([]pop.Individual, 100)
	for i := 0; i < 100; i++ {
		p.Members[i] = NewEnvInd(len(r.nodes)*2, rand.Float64()*400, r.Graph)
	}
	p.FitnessTests = 1
	p.Fitnesses = make([]int, 100)
	p.Elites = 4
	p.Size = 100
	p.Selection = selection.DeterministicTournament{
		TournamentSize:   2,
		ParentProportion: 2,
	}
	p.Pairing = pairing.Random{}
	p.GoalFitness = 1
	for i := 0; i < 200; i++ {
		if p.NextGeneration() {
			break
		}
	}

	rs := make([]render.Renderable, len(r.nodes))
	for i, n := range r.nodes {
		if rn, ok := n.(RenderNode); ok {
			rs[i] = rn.GetR()
		}
		// else ?
	}

	positions := make([]floatgeom.Point2, len(r.nodes))
	topMem, fitness := p.BestMember()
	fmt.Println("Top member fitness", fitness)
	topEnv := topMem.(*EnvInd)

	minX := 0.0
	minY := 0.0

	for i := 0; i < topEnv.F.Len(); i += 2 {
		pt := floatgeom.Point2{topEnv.Get(i), topEnv.Get(i + 1)}
		fmt.Println("Renderable position", pt.X(), pt.Y())
		if pt.X() < minX {
			minX = pt.X()
		}
		if pt.Y() < minY {
			minY = pt.Y()
		}
		positions[i/2] = pt
	}
	minPt := floatgeom.Point2{0, 0}
	if minX < 0 {
		minPt.Add(floatgeom.Point2{-minX, 0})
	}
	if minY < 0 {
		minPt.Add(floatgeom.Point2{0, -minY})
	}
	if minX < 0 || minY < 0 {
		for i, pt := range positions {
			positions[i] = pt.Add(minPt)
		}
	}

	for i, r := range rs {
		pos := positions[i]
		r.SetPos(pos.X(), pos.Y())
		fmt.Println("Renderable position", pos.X(), pos.Y())
	}

	r.CompositeR = render.NewCompositeR(rs)
}
