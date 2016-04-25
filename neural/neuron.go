package neural

import (
	"bytes"
	"strconv"
)

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
