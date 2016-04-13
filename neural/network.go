package neural

type Neuron interface {
	String() string
}

type Network interface {
	Get(x, y int) Neuron
	Slice(start, end int) Network
	SliceToEnd(start int) Network
	SliceFromStart(end int) Network
	Length() int
	Append(data interface{}) Network
	Make() Network
	Print()
	Fitness(inputs, expected [][]float64) int
}
