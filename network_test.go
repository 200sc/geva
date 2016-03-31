package neural

import (
	"testing"
)

func TestNetworkGeneration(*testing.T) {


	wOpt := FloatMutationOptions{
		0.33,
		0.05,
		5,
	}

	tOpt := FloatMutationOptions{
		0.33,
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
		0.05,
		0.02,
		0.10,
		0.10,
		0.10,
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

	network := GenerateNetwork(&nngOpt)

	network.print()
}