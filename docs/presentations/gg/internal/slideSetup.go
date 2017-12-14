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
		how,
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

func CodeSlide(ss *static.Slide, header string, list ...string) {
	ss.Append(show.Header(header))
	ss.Append(show.TxtSetFrom(Gnuolane44, .25, .35, 0, .04, list...)...)
}

func DoubleListSlide(ss *static.Slide, header string, list ...string) {
	ss.Append(show.Header(header))
	left := list[:len(list)/2]
	right := list[len(list)/2:]
	ss.Append(show.TxtSetFrom(Gnuolane44, .15, .35, 0, .07, left...)...)
	ss.Append(show.TxtSetFrom(Gnuolane44, .55, .35, 0, .07, right...)...)
}
