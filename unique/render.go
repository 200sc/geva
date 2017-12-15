package unique

import (
	"image"
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
	img   *image.RGBA
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

		// make r.img as large as need be

		// first increase all positions so all nodes have their leftmost point
		// greater than 0
		// then find the furthest out node positions (x+w,y+h) and set that as
		// the max bounds.

		// draw to r.img each node centered at their position in the population

		r.clean = true
	}
	// draw r.img to the buffer offset by xOff and yOff
}
