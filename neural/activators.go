package neural

import (
	"fmt"
	"math"
)

var (
	activatorTests = []float64{
		-4.0, -3.75, -3.5, -3.25,
		-3.0, -2.75, -2.5, -2.25,
		-2.0, -1.75, -1.5, -1.25,
		-1.0, -0.75, -0.5, -0.25,
		0.0, 0.75, 0.5, 0.25,
		1.0, 1.75, 1.5, 1.25,
		2.0, 2.75, 2.5, 2.25,
		3.0, 3.25, 3.5, 3.75,
		4.0,
	}
	testSize = 33
)

// This list egregiously copied from Wikipedia.
// To be expanded upon.

func Identity(x float64) float64 {
	return x
}

func BentIdentity(x float64) float64 {
	return ((math.Sqrt((math.Pow(x, 2) + 1)) - 1) / 2) + x
}

func Perceptron_Threshold(t float64) ActivatorFunc {
	return func(x float64) float64 {
		if x > t {
			return 1.0
		}
		return 0.0
	}
}

func Rectifier(x float64) float64 {
	return math.Max(x, 0.0)
}

func Rectifier_Parametric(a float64) ActivatorFunc {
	return func(x float64) float64 {
		if x >= 0 {
			return x
		}
		return a * x
	}
}

func Rectifier_Exponential(a float64) ActivatorFunc {
	return func(x float64) float64 {
		if x >= 0 {
			return x
		}
		return a * (math.Pow(math.E, x) - 1)
	}
}

func Softplus(x float64) float64 {
	return math.Log(1 + math.Pow(math.E, x))
}

func Softstep(x float64) float64 {
	return 1 / (1 + math.Pow(math.E, -1*x))
}

func Softsign(x float64) float64 {
	return x / (1 + math.Abs(x))
}

func Sinc(x float64) float64 {
	if x == 0 {
		return 1
	}
	return math.Sin(x) / x
}

// Bug(patrick)
// This currently results in negative fitnesses, which shouldn't be possible
// func Gaussian(x float64) float64 {
// 	return math.Pow(math.E, math.Pow(-1*x, 2))
// }

// For TanH, ArcTan, Sin
// just dump in math.Tanh / math.Arctan / math.Sin

/**
 * Print the activator function as an ASCII graph
 * Uses a staticly defined range
 */
func GraphPrintActivator(a ActivatorFunc) {
	m := make(map[float64]float64)
	for _, v := range activatorTests {
		m[v] = a(v)
	}
	grid := make([][]bool, testSize)
	for i := 0; i < testSize; i++ {
		grid[i] = make([]bool, testSize)
	}
	maxRow := 0
	minRow := 16
	for k, v := range m {
		gridCol := int((k + 4) * 4)
		gridRow := round((((v - 4) * 4) * -1))
		if gridRow < 0 || gridRow > testSize-1 {
			continue
		}
		grid[gridRow][gridCol] = true
		if gridRow > maxRow {
			maxRow = gridRow
		} else if gridRow < minRow {
			minRow = gridRow
		}
	}
	for y, row := range grid {
		if y > maxRow {
			continue
		} else if y < minRow {
			continue
		}
		for x, v := range row {
			if v {
				fmt.Print("■")
			} else {
				if x == 16 {
					fmt.Print("|")
				} else if y == 16 {
					fmt.Print("-")
				} else {
					fmt.Print(" ")
				}
			}
		}
		fmt.Println()
	}
}

func round(f float64) int {
	if f < -0.5 {
		return int(f - 0.5)
	}
	if f > 0.5 {
		return int(f + 0.5)
	}
	return 0
}
