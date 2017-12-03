package internal

import (
	"github.com/oakmound/oak/examples/slide/show"
	"github.com/oakmound/oak/examples/slide/show/static"
)

var (
	intro = SlideSetup{
		addIntro,
		2,
	}
)

func addIntro(i int, sslides []*static.Slide) {

	sslides[i].Append(
		show.Title("Generating Unique Game Mechanics through Co-Evolution"),
		show.TxtAt(Gnuolane44, "Patrick Stephen", .5, .5),
	)
	ListSlide(sslides[i+1], "Overview",
		"- Why",
		"- What",
		"- How",
	)
}
