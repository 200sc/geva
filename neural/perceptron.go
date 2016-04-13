package neural

import (
	"bytes"
	"strconv"
)

// A Perceptron has a list of places to send to
// and a mapping of places it receives from to weights.
// These lists are represented as integers, as a Perceptron has some
// presence in a "column" of Perceptrons-- it recieves
// from the previous column and sends to the following.
//
// For all Perceptrons which are not in an end column,
// it's assumed that they have at least one value in their
// Outputs and in their weights,
// for the purpose of mutation algorithms.
type Perceptron struct {
	threshold float64
	// For a performance boost and complexity reduction,
	// this could be replaced with a data structure of
	// a map which externally keeps track of an array of
	// keys for random element access
	Outputs map[int]bool
	weights map[int]float64
}

func (n Perceptron) String() string {
	var buffer bytes.Buffer

	buffer.WriteString("[")
	buffer.WriteString("t:")
	buffer.WriteString(strconv.FormatFloat(n.threshold, 'f', 2, 64))
	for k, v := range n.weights {
		buffer.WriteString("(")
		buffer.WriteString(strconv.Itoa(k))
		buffer.WriteString(",")
		buffer.WriteString(strconv.FormatFloat(v, 'f', 2, 64))
		buffer.WriteString(")")
	}

	if len(n.Outputs) > 0 {
		buffer.WriteString("->")
		i := 0
		for k := range n.Outputs {
			buffer.WriteString(strconv.Itoa(k))
			if i < len(n.Outputs)-1 {
				buffer.WriteString(",")
			}
			i++
		}
	}

	buffer.WriteString("] ")
	return buffer.String()
}
