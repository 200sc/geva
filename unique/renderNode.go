package unique

import "github.com/oakmound/oak/render"

// A RenderNode can act as a Node and
// a screen Renderable
type RenderNode interface {
	Node
	render.Renderable
}
