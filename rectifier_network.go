package neural

import (
	"fmt"
	"bytes"
	"math/rand"
	"strconv"
	"math"
)

// A neuron has a list of places to send to
// and a mapping of places it receives from to weights.
// These lists are represented as integers, as a neuron has some
// presence in a "column" of neurons-- it recieves
// from the previous column and sends to the following.
//
// For all Neurons which are not in an end column,
// it's assumed that they have at least one value in their
// outputs and in their weights, 
// for the purpose of mutation algorithms.
type RectifierNeuron []float64

type RectifierNetworkOutput struct {
	value float64
	index int 
}

func (n_p *RectifierNeuron) String() string {
	var buffer bytes.Buffer

	n := *n_p

	buffer.WriteString("[")

	if len(n) > 0 {
		for i, k := range n {
			buffer.WriteString(strconv.FormatFloat(k,'f',2,64))
			if i < len(n) - 1 {
				buffer.WriteString(",")
			}
		}
	}

	buffer.WriteString("] ")
	return buffer.String()
}	

type RectifierNetwork [][]RectifierNeuron


/**
 * Take a network and duplicate it
 */
func (nn_p *RectifierNetwork) copy() RectifierNetwork {
	
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
	cOpt := *nnOpt.columnOptions

	nn := RectifierNetwork{}

	// Set up the input column
	inputColumn := make([]RectifierNeuron, nnOpt.inputs)

	nn = append(nn, inputColumn)

	// Set up the output column
	outputColumn := make([]RectifierNeuron, nnOpt.outputs)

	nn = append(nn, outputColumn)

	columnCount := rand.Intn(nnOpt.maxColumns - nnOpt.minColumns) + nnOpt.minColumns

	for i := 0; i < columnCount; i++ {
		nn = *(nn.addColumn(&cOpt))
	}

	nn.Print()

	for i := 0; i < nnOpt.baseMutations; i++ {
		nn = *(nn.Mutate(&nnOpt.RectifierNetworkMutationOptions))	
	}

	return &nn
}

// Todo: Improve this
func (nn_p *RectifierNetwork) Print() {
	for _,col := range *nn_p{
		for _,n := range col {
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
func (nn_p *RectifierNetwork) Run(inputs []float64) []float64 {

	nn := *nn_p

	doneCh := make(chan RectifierNetworkOutput)

	channels := make([][]chan float64, len(nn))
	for x,col := range nn {
		channels[x] = make([]chan float64, len(col))
	}
	for x, col := range nn {
		for y,neuron := range col {
			if x == len(nn) - 1 {
				go func(inputChannel chan float64, doneCh chan RectifierNetworkOutput,
					    inputLength int, y int) {
					
					out := 0.0

					for i := 0; i < inputLength; i++ {
						out += <-inputChannel
					}
					close(inputChannel)

					out = math.Max(out,0.0)
					
					doneCh <- RectifierNetworkOutput{out, y}

				}(channels[x][y], doneCh, len(channels[x-1]), y)
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
					// inputs for our value.
					for i := 0; i < inputLength; i++ {
						out += <-inputChannel
					}
					close(inputChannel)

					// This is the 'rectifier'.
					// Instead of sending a signal 0 or 1
					// based on a threshold, we send our actual
					// value (so long as it exceeds 0).
					out = math.Max(out,0.0)

					// As above, we apply the next column's
					// weights as we send them off.
					for i, weight := range n {
						channelColumn[i] <- out * weight
					}
				}(neuron, channels[x][y], channels[x+1], len(channels[x-1]))
			}
		}
	}

	// Send the first row their initial values
	for i, ch := range channels[0] {
		ch <- inputs[i]
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

// High fitness is bad, and vice versa.
func (n_p *RectifierNetwork) Fitness(inputs, expected [][]float64) int {
	fitness := 1.0
	for i := range inputs {
		output := (*n_p).Run(inputs[i])
		for j := range output {
			fitness += math.Abs(output[j] - expected[i][j])
		}
	}
	return int(math.Ceil(fitness))
}
