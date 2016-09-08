package population

type Individual interface {
	Fitness(input, expected [][]float64) int
	Mutate()
	Crossover(other Individual) Individual
	CanCrossover(other Individual) bool
	Print()
}
