package neural

import (
	"testing"
	"fmt"
)

func TestNetworkGeneration(*testing.T) {

	wOpt := FloatMutationOptions{
		0.66,
		0.05,
		5,
	}

	tOpt := FloatMutationOptions{
		0.66,
		0.05,
		5,
	}

	nmOpt := NeuronMutationOptions{
		tOpt,
		wOpt,
	}

	ngOpt := NeuronGenerationOptions{
		1,
		3,
		0.8,
		1,
	}

	cgOpt := ColumnGenerationOptions{
		2,
		16,
		&ngOpt,
	}

	nnmOpt := NetworkMutationOptions{
		&nmOpt,
		&cgOpt,
		0.02,
		0.05,
		0.03,
		0.03,
		0.05,
		0.01,
		0.02,
		0.10,
	}

	nngOpt := NetworkGenerationOptions{
		nnmOpt,
		10,
		20,
		3,
		4,
		50,
	}

	network := GenerateNetwork(&nngOpt)
	network.Print()
	for {	
		fmt.Println(network.Run([]bool{true,true,true}))
		fmt.Println(network.Run([]bool{true,true,false}))
		fmt.Println(network.Run([]bool{true,false,true}))
		fmt.Println(network.Run([]bool{true,false,false}))
		fmt.Println(network.Run([]bool{false,true,true}))	
		fmt.Println(network.Run([]bool{false,true,false}))
		fmt.Println(network.Run([]bool{false,false,true}))
		fmt.Println(network.Run([]bool{false,false,false}))
	}
}