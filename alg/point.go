package alg

import "math"

type Point struct {
	X, Y float64
}

func IntPoint(x, y int) Point {
	return Point{float64(x), float64(y)}
}

func Distance(a, b Point) float64 {
	return math.Sqrt(math.Pow(b.X-a.X, 2) + math.Pow(b.Y-a.Y, 2))
}
