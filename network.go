package neural

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
	// For a performance boost and complexity reduction,
	// this could be replaced with a data structure of 
	// a map which externally keeps track of an array of 
	// keys for random array element access
	outputs map[int]bool
	threshold float
	weights map[int]float
	val int
}

type Network struct {
	columns [][]Neuron
}

/**
 * Take a network and duplicate it
 */
func (nn *Network) copy() Network {

	var columns [][]Neuron
	var ic InputColumn

	for i := range nn.columns {
    	columns[i] = make([]Neuron, len(nn.columns[i]))
    	copy(columns[i], nn.columns[i])
	}

	for i := range nn.ic {
		ic[i] = make([]int len(nn.ic[i]))
		copy(ic[i], nn.ic[i])
	}
	
	return Network{ic, columns}
}

/**
 * Run some input through a neural network.
 * This returns the network's output column.
 */
func (nn *Network) run(inputs []bool) bool {

	channels = [][]chan int

	doneCh := make(chan bool)

	for x,col := range nn.columns {
		channels = append(channels, []chan int)
		for y,neuron := range col {

			channels[x] = append(channels[x], make(chan int))

			// Create a goroutine for every channel index
			// which accepts the neuron at that index
			go func(n_p *Neuron, inputChannel chan int, channelColumn []chan int, doneCh chan bool) {
				inputs := make(map[int]int)

				n := *n_p

				// Wait on all expected inputs.
				// An input is the form of a positive
				// or negative index from the previous column.
				// If it's negative, the signal we received from
				// that column is false. If it's positive, true.
				// 
				// Those signals are then mapped according to
				// their column, as our weights are also column-
				// indexed. 
				for i := range n.weights {
					input := <-inputChannel
					if input < 0 {
						// This 1 subraction is to
						// compensate for the 0-indexing
						// to 1-indexing shift
						inputs[(-1*input)-1] = 0
					} else {
						inputs[(input-1)] = 1
					}
					close(inputChannel)
				}

				// Sum the signals received
				// as according to our weights.
				out := 0
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
					result = -1 * (n.y + 1)
					// -1 represents false,
					// as 0 represents uninitialized. 
					n.val = -1
				} else {
					result = n.y + 1
					n.val = 1
				}

				// Send our result to all of 
				// this neuron's output channels.
				// If this is the last column,
				// outputs should be empty.
				for outY,_ := range n.outputs {
					channelColumn[outY] <- result
				}

				// Send a notification we finished if
				// we're an output neuron
				// TODO: ENSURE NO NEURON HAS A DEAD
				// END EXCEPT FOR OUTPUT NEURONS
				if len(n.outputs) == 0 {
					doneCh <- true
				}

			}(&neuron, channels[x][y], channels[x+1], doneCh)
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
	for i := 0; i < len(nn.columns[len(nn.columns)-1]); i++ {
		<-doneCh
	}
	close(doneCh)

	// Create an array of the last columns' values to return
	output = []int
	for _,neuron := range nn.columns[len(nn.columns)-1] {
		output = append(output, neuron.val)
	}

	return output
}

