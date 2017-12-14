package internal

import (
	"github.com/oakmound/oak/examples/slide/show"
	"github.com/oakmound/oak/examples/slide/show/static"
)

var (
	how = SlideSetup{
		addHow,
		8,
	}
)

func addHow(i int, sslides []*static.Slide) {

	sslides[i].Append(
		show.Title("A Game Mechanic Generator, and Evaluator"),
	)

	ListSlide(sslides[i+1], "The Developer Interface:",
		"- Can be queried for mechanics",
		"- Can be mutated",
		"- Can cross over with other developers",
	)

	DoubleListSlide(sslides[i+2], "A Developer Instance:",
		"- Full Environment Size Range",
		"- Initialized Variables Size Range",
		"- Goal Environment Size Range",
		"- Distance from Start to Goal Range",
		"- Starting Float Value Range",
		"- Action Count Range",
		"- Potential Types of Actions",
		"- Weighting on Action Type Selection",
		"- A Mutatable Mutation Method",
		"- A Mutatable Crossover Method",
	)

	ListSlide(sslides[i+3], "The Player Interface:",
		"- Can play mechanics and report how much they liked them",
		"- Play(Mechanic, maxPlayTime) Enjoyment",
		"- Enjoyment on a scale from 0 to 1",
	)

	ListSlide(sslides[i+4], "A Player Instance:",
		"- Has some fitness it considers winning",
		"- Has a desired amount of time before it wins",
		"- Evolves a series of actions to take to play the mechanic",
		"- If they win too easily, or can't win, reports lesser Enjoyment",
	)

	CodeSlide(sslides[i+5], "Algorithm",
		"Generate N Developers",
		"Loop M times:",
		"- Generate O Players",
		"- Evenly distribute Players to Devs",
		"- Loop P times:",
		"- - Create a Mechanic for each Dev",
		"- - Have each Player play their Dev's Mechanic",
		"- - If Players don't enjoy the game (rand() > Enjoyment), move them",
		"- Sort Devs by remaining Players",
		"- Evolve Dev population (Select, Crossover, Mutate)",
	)

	ListSlide(sslides[i+6], "Results So Far",
		"Playerbase is biased towards easier games.",
		"Developers create games with many (> 100) actions,",
		"but where the initial state ~= the goal state.",
		"So long as players don't choose the rare bad actions, they win.",
	)

	ListSlide(sslides[i+7], "Remaining Work",
		"- More Evolutionary Layers",
		"- Uniqueness Analysis and Graphing",
		"- Applying Results to a Game",
		"- Additional variety in mechanics (continuous, awareness of environment)",
	)
}
