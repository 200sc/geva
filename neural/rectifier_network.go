package neural

import (
	"fmt"
	"math"
	"math/rand"
)

// A neuron has a list of places to send to
// and a mapping of places it receives from to weights.
// These lists are represented as integers, as a neuron has some
// presence in a "column" of neurons-- it recieves
// from the previous column and sends to the following.
//
// For all Neurons which are not in an end column,
// it's assumed that they have at least one value in their
// Outputs and in their weights,
// for the purpose of mutation algorithms.

type RectifierNetworkOutput struct {
	value float64
	index int
}

type RectifierNetwork [][]RectifierNeuron

/**
 * Take a network and duplicate it
 */
func (nn_p *RectifierNetwork) Copy() RectifierNetwork {

	nn := *nn_p

	var newNetwork RectifierNetwork

	for i := range nn {
		newNetwork = append(newNetwork, make([]RectifierNeuron, len(nn[i])))
		copy(newNetwork[i], nn[i])
	}

	return newNetwork
}

func GenerateRectifierNetwork(nOpt_p *RectifierNetworkGenerationOptions) *RectifierNetwork {

	nnOpt := *nOpt_p
	cOpt := *nnOpt.ColumnOptions

	nn := RectifierNetwork{}

	// Set up the input column
	inputColumn := make([]RectifierNeuron, nnOpt.Inputs)

	nn = append(nn, inputColumn)

	// Set up the output column
	outputColumn := make([]RectifierNeuron, nnOpt.Outputs)

	nn = append(nn, outputColumn)

	// reset the input column to give it axons
	for i := 0; i < len(inputColumn); i++ {
		nn[0][i] = make(RectifierNeuron, len(outputColumn))
		nn.replaceNeuron(0, i, cOpt.DefaultAxonWeight)
	}

	columnCount := rand.Intn(nnOpt.MaxColumns-nnOpt.MinColumns) + nnOpt.MinColumns

	for i := 0; i < columnCount; i++ {
		nn = *(nn.addColumn(&cOpt))
	}

	for i := 0; i < nnOpt.BaseMutations; i++ {
		nn = *(nn.Mutate(&nnOpt.RectifierNetworkMutationOptions))
	}

	return &nn
}

// Todo: Improve this
func (nn RectifierNetwork) Print() {
	for _, col := range nn {
		for _, n := range col {
			fmt.Print(n.String())
		}
		fmt.Println("")
	}
	fmt.Println("")
}

/**
 * Run some input through a neural network.
 * This returns the network's output column.
 */
func (nn_p *RectifierNetwork) Run(Inputs []float64) []float64 {

	nn := *nn_p

	doneCh := make(chan RectifierNetworkOutput)

	channels := make([][]chan float64, len(nn))
	for x, col := range nn {
		channels = append(channels, []chan float64{})
		for range col {
			channels[x] = append(channels[x], make(chan float64))
		}
	}
	for x, col := range nn {
		for y, neuron := range col {
			var l int
			if x == 0 {
				l = 1
			} else {
				l = len(channels[x-1])
			}
			if x == len(nn)-1 {
				go func(inputChannel chan float64, doneCh chan RectifierNetworkOutput,
					inputLength int, y int) {

					out := 0.0

					for i := 0; i < inputLength; i++ {
						out += <-inputChannel
					}
					close(inputChannel)

					out = math.Max(out, 0.0)

					doneCh <- RectifierNetworkOutput{out, y}

				}(channels[x][y], doneCh, l, y)
			} else {

				go func(n RectifierNeuron, inputChannel chan float64,
					channelColumn []chan float64, inputLength int) {

					out := 0.0

					// Compared to network.go,
					// this system has each neuron
					// know the weights of the following layer,
					// instead of knowing where they output to.
					// It also has each neuron send to every
					// neuron in the following layer-- a 0
					// weight is equivalent to not accepting
					// some input.
					//
					// At this stage, that means we can
					// just sum all of our already-weighted
					// Inputs for our value.
					for i := 0; i < inputLength; i++ {
						out += <-inputChannel
					}
					close(inputChannel)

					// This is the 'rectifier'.
					// Instead of sending a signal 0 or 1
					// based on a threshold, we send our actual
					// value (so long as it exceeds 0).
					out = math.Max(out, 0.0)

					// As above, we apply the next column's
					// weights as we send them off.
					for i, weight := range n {
						channelColumn[i] <- out * weight
					}
				}(neuron, channels[x][y], channels[x+1], l)
			}
		}
	}

	// Send the first row their initial values
	for i, ch := range channels[0] {
		ch <- Inputs[i]
	}

	// We need to wait here, on the last columns being populated
	output := make([]float64, len(nn[len(nn)-1]))
	for i := 0; i < len(nn[len(nn)-1]); i++ {
		recieved := <-doneCh
		output[recieved.index] = recieved.value
	}
	close(doneCh)

	return output
}

func (n RectifierNetwork) Get(x, y int) Neuron {
	return n[x][y]
}
func (n RectifierNetwork) Slice(start, end int) Network {
	return n[start:end]
}
func (n RectifierNetwork) SliceToEnd(start int) Network {
	return n[start:]
}
func (n RectifierNetwork) SliceFromStart(end int) Network {
	return n[:end]
}
func (n RectifierNetwork) Length() int {
	return len(n)
}
func (n RectifierNetwork) Append(data interface{}) Network {
	n = append(n, data.(RectifierNetwork)...)
	return n
}
func (n RectifierNetwork) Make() Network {
	out := make(RectifierNetwork, 0)
	return out
}

// High fitness is bad, and vice versa.
func (n RectifierNetwork) Fitness(Inputs, expected [][]float64) int {
	fitness := 1.0
	for i := range Inputs {
		output := n.Run(Inputs[i])
		for j := range output {
			fitness += math.Abs(output[j] - expected[i][j])
		}
	}
	return int(math.Ceil(fitness))
}
