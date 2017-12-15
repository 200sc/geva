package unique

import (
	"image/draw"
	"math/rand"

	"github.com/200sc/geva/pop"
	"github.com/200sc/geva/selection"

	"github.com/oakmound/oak/render"
)

var (
	_ render.Renderable = &Render{}
)

type Render struct {
	*Graph
	render.LayeredPoint
	clean bool
}

func (r *Render) SetGraph(g *Graph) {
	r.Graph = g
	r.clean = false
}

func (r *Render) Draw(buff draw.Image) {
	r.DrawOffset(buff, 0, 0)
}
func (r *Render) DrawOffset(buff draw.Image, xOff, yOff float64) {
	if !r.clean {
		// there's necessarily going to be error if the nodes
		// represent more than two variables. We move elements to
		// minimize error, up to some number of tries.
		//
		// Is this a genetic algorithm problem
		// yes
		p := pop.Population{}
		p.Members = make([]pop.Individual, 100)
		for i := 0; i < 100; i++ {
			p.Members[i] = NewEnvInd(len(r.nodes)*2, rand.Float64()*100, r.Graph)
		}
		p.FitnessTests = 1
		p.Elites = 4
		p.Size = 100
		p.Selection = selection.DeterministicTournament{
			TournamentSize:   2,
			ParentProportion: 2,
		}
		p.GoalFitness = 1
		for i := 0; i < 50; i++ {
			if p.NextGeneration() {
				break
			}
		}

		r.clean = true
	}
}
