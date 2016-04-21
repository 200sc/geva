package neural

type Neuron interface {
	String() string
}

type Network interface {
	Get(x, y int) Neuron
	Set(x, y int, i interface{})
	Slice(start, end int) Network
	SliceToEnd(start int) Network
	SliceFromStart(end int) Network
	Length() int
	ColLength(i int) int
	Append(data interface{}) Network
	Make(size int) Network
	CopyStructure() Network
	Print()
	Fitness(inputs, expected [][]float64) int
}
