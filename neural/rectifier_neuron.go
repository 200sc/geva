package neural

import (
	"bytes"
	"strconv"
)

type RectifierNeuron []float64

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