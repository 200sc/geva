package eda

import (
	"math"

	"github.com/200sc/go-dist/floatrange"
)

// BoxMueller is shorthand for BoxMuellerSin
func BoxMueller(fr floatrange.Range) float64 {
	return BoxMuellerSin(fr)
}

// BoxMuellerSin returns Rsin(Theta)
func BoxMuellerSin(fr floatrange.Range) float64 {
	return math.Sqrt(-2*math.Log(fr.Poll())) * math.Sin(2*math.Pi*fr.Poll())
}

// BoxMuellerCos returns Rcos(Theta)
func BoxMuellerCos(fr floatrange.Range) float64 {
	return math.Sqrt(-2*math.Log(fr.Poll())) * math.Cos(2*math.Pi*fr.Poll())
}
