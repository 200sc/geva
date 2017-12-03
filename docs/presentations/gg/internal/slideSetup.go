package internal

import (
	"github.com/oakmound/oak/examples/slide/show"
	"github.com/oakmound/oak/examples/slide/show/static"
)

var (
	Setups = []SlideSetup{
		intro,
		why,
		what,
	}
)

type SlideSetup struct {
	Add func(int, []*static.Slide)
	Len int
}

func ListSlide(ss *static.Slide, header string, list ...string) {
	ss.Append(show.Header(header))
	ss.Append(show.TxtSetFrom(Gnuolane44, .25, .35, 0, .07, list...)...)
}
