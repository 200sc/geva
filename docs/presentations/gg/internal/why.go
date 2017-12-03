package internal

import (
	"github.com/oakmound/oak/examples/slide/show"
	"github.com/oakmound/oak/examples/slide/show/static"
)

var (
	why = SlideSetup{
		addWhy,
		4,
	}
)

func addWhy(i int, sslides []*static.Slide) {
	sslides[i].Append(
		show.Title("Goal: An Agent that Generates Video Games"),
	)
	sslides[i+1].Append(
		show.Title("How Can We Generate Video Games?"),
	)
	ListSlide(sslides[i+2], "What Is A Video Game?",
		"- Art",
		"- SFX",
		"- Music",
		"- Story",
		"- Voice Acting",
		"- Gameplay",
	)
	ListSlide(sslides[i+3], "What Is A Video Game?",
		"- Art (NNs)",
		"- SFX (NNs)",
		"- Music (NNs)",
		"- Story (NNs)",
		"- Voice Acting (NNs)",
		"- Gameplay (?)",
	)
}
