package neural

import (
	"testing"
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
		0.5,
		1,
	}

	cgOpt := ColumnGenerationOptions{
		2,
		8,
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
		4,
		10,
		3,
		1,
		20,
	}

	for i := 0; i < 100000; i++ {
		GenerateNetwork(&nngOpt)
		//network.print()
	}
}