package neural

import (
	"math"
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

func Gaussian(x float64) float64 {
	return math.Pow(math.E, math.Pow(-1*x, 2))
}

// For TanH, ArcTan, Sin
// just dump in math.Tanh / math.Arctan / math.Sin
