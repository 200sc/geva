package unique

import (
	"testing"

	"github.com/200sc/go-dist/floatrange"
	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/stretchr/testify/assert"
)

func TestScale(t *testing.T) {
	rng := floatrange.NewLinear(-100000, 1000000)
	in := make([]floatgeom.Point2, 10000)
	for i := 0; i < 10000; i++ {
		in[i] = floatgeom.Point2{rng.Poll(), rng.Poll()}
	}
	r := floatgeom.NewRect2(30, 30, 640, 320)
	positions := scaleToRect(r, in...)
	for _, p := range positions {
		assert.True(t, p.X() >= 30 && p.Y() >= 30 && p.X() <= 640 && p.Y() <= 320)
	}
}
