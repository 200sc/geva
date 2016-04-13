package neural

import (
	"fmt"
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

	nmOpt := PerceptronMutationOptions{
		tOpt,
		wOpt,
	}

	ngOpt := PerceptronGenerationOptions{
		1,
		3,
		0.8,
		1,
	}

	cgOpt := PerceptronColumnGenerationOptions{
		2,
		16,
		&ngOpt,
	}

	nnmOpt := PerceptronNetworkMutationOptions{
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

	nngOpt := PerceptronNetworkGenerationOptions{
		nnmOpt,
		10,
		20,
		3,
		4,
		50,
	}

	network := GeneratePerceptronNetwork(&nngOpt)
	network.Print()
	fmt.Println(network.Run([]bool{true, true, true}))
	fmt.Println(network.Run([]bool{true, true, false}))
	fmt.Println(network.Run([]bool{true, false, true}))
	fmt.Println(network.Run([]bool{true, false, false}))
	fmt.Println(network.Run([]bool{false, true, true}))
	fmt.Println(network.Run([]bool{false, true, false}))
	fmt.Println(network.Run([]bool{false, false, true}))
	fmt.Println(network.Run([]bool{false, false, false}))
}

func TestRectifierNetworkGeneration(t *testing.T) {
	wOpt := FloatMutationOptions{
		0.40,
		0.20,
		5,
	}

	cgOpt := RectifierColumnGenerationOptions{
		2,
		16,
		0.1,
	}

	nnmOpt := RectifierNetworkMutationOptions{
		&wOpt,
		&cgOpt,
		0.02,
		0.06,
		0.06,
		0.01,
		0.01,
		0.33,
	}

	nngOpt := RectifierNetworkGenerationOptions{
		nnmOpt,
		5,
		20,
		3,
		4,
		50,
	}

	network := GenerateRectifierNetwork(&nngOpt)
	network.Print()
	fmt.Println(network.Run([]float64{1.0, 1.0, 1.0}))
	fmt.Println(network.Run([]float64{1.0, 1.0, -1.0}))
	fmt.Println(network.Run([]float64{1.0, -1.0, 1.0}))
	fmt.Println(network.Run([]float64{1.0, -1.0, -1.0}))
	fmt.Println(network.Run([]float64{-1.0, 1.0, 1.0}))
	fmt.Println(network.Run([]float64{-1.0, 1.0, -1.0}))
	fmt.Println(network.Run([]float64{-1.0, -1.0, 1.0}))
	fmt.Println(network.Run([]float64{-1.0, -1.0, -1.0}))
}

func BenchmarkNetworkGeneration(b *testing.B) {
	for i := 0; i < b.N; i++ {
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

		network := GeneratePerceptronNetwork(&nngOpt)
		network.Run([]bool{true, true, true})
	}
}

func BenchmarkRectifierNetworkGeneration(b *testing.B) {
	for i := 0; i < b.N; i++ {
		wOpt := FloatMutationOptions{
			0.66,
			0.05,
			5,
		}

		cgOpt := RectifierColumnGenerationOptions{
			2,
			16,
			0.1,
		}

		nnmOpt := RectifierNetworkMutationOptions{
			&wOpt,
			&cgOpt,
			0.02,
			0.05,
			0.05,
			0.01,
			0.02,
			0.10,
		}

		nngOpt := RectifierNetworkGenerationOptions{
			nnmOpt,
			10,
			20,
			3,
			4,
			50,
		}

		network := GenerateRectifierNetwork(&nngOpt)
		network.Run([]float64{1.0, 1.0, 1.0})
	}
}
