package neural

import (
	"fmt"
	"bytes"
	"math/rand"
	"strconv"
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
type Neuron struct {
	threshold float64
	// For a performance boost and complexity reduction,
	// this could be replaced with a data structure of 
	// a map which externally keeps track of an array of 
	// keys for random element access
	outputs map[int]bool
	weights map[int]float64
}

func (n_p *Neuron) String() string {
	var buffer bytes.Buffer

	n := *n_p

	buffer.WriteString("[")
	buffer.WriteString("t:")
	buffer.WriteString(strconv.FormatFloat(n.threshold,'f',2,64))
	for k, v := range n.weights {
		buffer.WriteString("(")
		buffer.WriteString(strconv.Itoa(k))
		buffer.WriteString(",")
		buffer.WriteString(strconv.FormatFloat(v,'f',2,64))
		buffer.WriteString(")")
	}

	if len(n.outputs) > 0 {
		buffer.WriteString("->")
		i := 0
		for k := range n.outputs {
			buffer.WriteString(strconv.Itoa(k))
			if i < len(n.outputs) - 1 {
				buffer.WriteString(",")
			}
			i++
		}
	}

	buffer.WriteString("] ")
	return buffer.String()
}	

type Network [][]Neuron


/**
 * Take a network and duplicate it
 */
func (nn_p *Network) copy() Network {
	
	nn := *nn_p

	var newNetwork Network

	for i := range nn {
    	newNetwork = append(newNetwork, make([]Neuron, len(nn[i])))
    	copy(newNetwork[i], nn[i])
	}
	
	return newNetwork
}

func GenerateNetwork(nOpt_p *NetworkGenerationOptions) *Network {

	nnOpt := *nOpt_p
	cOpt := *nnOpt.columnOptions
	nOpt := *cOpt.neuronOptions

	nn := Network{}

	// Set up the input column
	inputColumn := []Neuron{}

	for i := 0; i < nnOpt.inputs; i++ {
		n := Neuron{threshold: nOpt.defaultThreshold,
					weights: map[int]float64{0:nOpt.defaultAxonWeight},
					outputs: make(map[int]bool),
				}
		inputColumn = append(inputColumn, n)
	}

	nn = append(nn, inputColumn)

	// Set up the output column
	outputColumn := []Neuron{}

	for i := 0; i < nnOpt.outputs; i++ {
		n := Neuron{
			threshold: nOpt.defaultThreshold,
			weights: make(map[int]float64),
			outputs: make(map[int]bool),
		}
		outputColumn = append(outputColumn, n)
	}

	nn = append(nn, outputColumn)

	columnCount := rand.Intn(nnOpt.maxColumns - nnOpt.minColumns) + nnOpt.minColumns

	for i := 0; i < columnCount; i++ {
		nn = *(nn.addColumn(&cOpt))
	}

	for i := 0; i < nnOpt.baseMutations; i++ {
		nn = *(nn.Mutate(&nnOpt.NetworkMutationOptions))	
	}

	return &nn
}

// Todo: Improve this
func (nn_p *Network) Print() {
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
func (nn_p *Network) Run(inputs []bool) []bool {

	nn := *nn_p

	doneCh := make(chan int)

	channels := make([][]chan int, len(nn))
	for x,col := range nn {
		channels[x] = make([]chan int, len(col))
	}
	for x, col := range nn {
		for y,neuron := range col {
			if x == len(nn) - 1 {
				go func(n Neuron, inputChannel chan int, doneCh chan int, y int) {
					inputs := make(map[int]float64)

					for _ = range n.weights {
						input := <-inputChannel
						if input < 0 {
							inputs[(-1*input)-1] = 0
						} else {
							inputs[(input-1)] = 1
						}
					}
					close(inputChannel)

					out := 0.0
					for i,weight := range n.weights {
						out += weight * inputs[i]
					}
					
					var result int 
					if out <= n.threshold {
						// If we didn't add one, the 0-index
						// would be indistinguishable from others 
						result = -1 * (y + 1)
					} else {
						result = y + 1
					}

					doneCh <- result

				}(neuron, channels[x][y], doneCh, y)
			} else {

				// Create a goroutine for every channel index
				// which accepts the neuron at that index
				go func(n Neuron, inputChannel chan int, channelColumn []chan int, y int, x int) {
					inputs := make(map[int]float64)

					// Wait on all expected inputs.
					// An input is the form of a positive
					// or negative index from the previous column.
					// If it's negative, the signal we received from
					// that column is false. If it's positive, true.
					// 
					// Those signals are then mapped according to
					// their column, as our weights are also column-
					// indexed. 
					for i := 0; i < len(n.weights); i++ {
						input := <-inputChannel
						if input < 0 {
							// This 1 subraction is to
							// compensate for the 0-indexing
							// to 1-indexing shift
							inputs[(-1*input)-1] = 0
						} else {
							inputs[(input-1)] = 1
						}
					}
					// Technically we probably
					// don't need to close these channels
					// but it's best practice and, perhaps
					// more importantly, lets us know if
					// we're majorly screwing up by sending
					// to a channel we're supposed to be
					// done with. 
					close(inputChannel)

					// Sum the signals received
					// as according to our weights.
					out := 0.0
					for i,weight := range n.weights {
						out += weight * inputs[i]
					}
					
					// Our result is either negative or 
					// postive based on whether our sum
					// exceeds our threshold.
					var result int 
					if out <= n.threshold {
						// If we didn't add one, the 0-index
						// would be indistinguishable from others 
						result = -1 * (y + 1)
					} else {
						result = y + 1
					}

					// Send our result to all of 
					// this neuron's output channels.
					// If this is the last column,
					// outputs should be empty.
					for outY,_ := range n.outputs {
						channelColumn[outY] <- result
					}

				}(neuron, channels[x][y], channels[x+1], y, x)
			}
		}
	}

	// Send the first row their initial values
	for i, ch := range channels[0] {
		if inputs[i] {
			ch <- 1
 		} else {
 			ch <- -1
 		}
	}

	// We need to wait here, on the last columns being populated
	output := make([]bool, len(nn[len(nn)-1]))
	for i := 0; i < len(nn[len(nn)-1]); i++ {
		val := <-doneCh
		if val < 0 {
			// This 1 subraction is to
			// compensate for the 0-indexing
			// to 1-indexing shift
			output[(-1*val)-1] = false
		} else {
			output[val-1] = true
		}
	}
	close(doneCh)

	return output
}

// High fitness is bad, and vice versa.
func (*n_p Network) Fitness(inputs, expected [][]bool) int {
	fitness := 1
	for i := range inputs {
		output := (*n_p).Run(inputs[i])
		for j := range output {
			if output[j] != expected[i][j] {
				fitness++
			}
		}
	}
	return fitness
}
