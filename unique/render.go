package unique

import (
	"image/draw"

	"github.com/oakmound/oak/alg/floatgeom"

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

		positions := make([]floatgeom.Point2, len(r.nodes))

		// there's necessarily going to be error if the nodes
		// represent more than two variables. We move elements to
		// minimize error, up to some number of tries.
		//
		// Is this a genetic algorithm problem
		// yes

	}
}
