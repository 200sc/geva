package unique

type Option func(*Graph)

func MinDistance(i float64) Option {
	return func(g *Graph) {
		g.MinDistance = i
	}
}
