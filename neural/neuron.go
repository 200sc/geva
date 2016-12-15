package neural

import (
	"bytes"
	"strconv"
)

// A Neuron is a list of weights.
// Clasically, the weights on a neuron would normally
// represent what that neuron would multiply its inputs
// by to obtain it's value.
//
// These weights do not represent that. These weights
// represent what this neuron should multiply its input
// by before sending it to the next column, for each
// element in the next column.
//
// Effectively, each neuron receives pre-weighted values.
// There's no difference in how the neurons function--
// interpret a neuron's weights as the set of weights
// from the previous column where the index in each
// previous column's neuron's weights matches the index
// of the desired neuron in the following column,
// if you so choose.
//
// All Neurons connect to all Neurons in the following column.
// A weight of 0.0 represents what would classically be no
// connection.
//
// There probably isn't a significant difference in performance between
// these two representations. The significant implementation difference
// is where the delay happens on channel sending-- does it happen
// as signals are sent, or does it happen as they are received?
type Neuron []float64

/**
 * Obtain a string represenation of a neuron
 **/
func (n Neuron) String() string {
	var buffer bytes.Buffer

	buffer.WriteString("[")

	if len(n) > 0 {
		for i, k := range n {
			buffer.WriteString(strconv.FormatFloat(k, 'f', 2, 64))
			if i < len(n)-1 {
				buffer.WriteString(",")
			}
		}
	}

	buffer.WriteString("] ")
	return buffer.String()
}
