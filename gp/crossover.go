package gp

type GPCrossover interface {
	Crossover(a, b *GP) *GP
}
