package player

import (
	"github.com/200sc/geva/gg/dev"
)

type Player interface {
	Play(*dev.Mechanic, int) float64
}
