package internal

import (
	"github.com/oakmound/oak/examples/slide/show"
	"github.com/oakmound/oak/examples/slide/show/static"
)

var (
	what = SlideSetup{
		addWhat,
		3,
	}
)

func addWhat(i int, sslides []*static.Slide) {

	sslides[i].Append(
		show.Title("A Generatable Game Mechanic"),
	)

	ListSlide(sslides[i+1], "What makes up a Mechanic?",
		"- Starting Environment",
		"- Player Actions",
		"- Goal Environment(s)",
	)

	ListSlide(sslides[i+2], "What makes up a Mechanic?",
		"- An initial vector of floats",
		"- Functions that modify particular vector values",
		"- A boolean check that a vector of floats is a goal state",
	)
}
