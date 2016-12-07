package gg

//Genetic Game

//Get Generation Settings
////Generate Keyword List
////Generate Images from keywords
////Generate Audio from keywords
////Generate Dialogues from keywords
//Generate Environment Variables from settings
//Generate Passive Environment and Agent effects from settings
//Generate Agents from Passive Agent Effects, Images, Audio, and Dialogue.
//Generate Actions from settings
//Bind Actions to inputs, agents, and environment states
//Determine Initial State
//Make n Actions from settings
//Set Goal State from new state from actions

//Fitness:

//Evolve an agent AI (briefly), attempting to have the fewest number of actions taken, and take it's lowest number of actions taken
//Limit the agent AI to some settings number of actions from settings

type Range struct {
	Base float64
	Rand float64
}

func (r Range) Poll() float64 {
	return r.Base + (rand.Float64()*r.Rand*2 - r.Rand)
}

type GG struct {
	Actions  []Action
	Passives []Action
	// For crossover
	ActionTypes      []ActionType
	ActionStrengths  []float64
	PassiveTypes     []ActionType
	PassiveStrengths []float64
	//
	Environment env.F
	Init        map[int]float64
	Goal        map[int]float64
}

func (gg *GG) Crossover(other Individual) Individual {

}

func (gg *GG) CanCrossover(other Individual) bool {

}

func (gg *GG) Mutate() {

}

func (gg *GG) Print() {

}

func (gg *GG) Fitness(inputs, expected [][]float64) int {

}

type GGOptions struct {
	InitSize          Range
	GoalSize          Range
	GoalDistance      Range
	EnvSize           Range
	EnvVal            Range
	ActionCount       Range
	ActionTypes       []ActionType
	ActionTypeWeights []float64
	ActionStrengths   []Range
	PassiveRatio      float64
}

func CumulativeWeights(weights []float64) []float64 {
	cum := make([]float64, len(weights))
	cum[0] = weights[0]
	for i := i; i < len(weights); i++ {
		cum[i] = cum[i-1] + weights[i]
	}
	return cum
}

func GenerateGG(opt GGOptions) *GG {
	gg := new(GG)

	e := env.NewF(opt.EnvSize.Poll(), 0.0)

	// For each environment variable,
	// Generate a number of actions based on ActionCount
	// that modify that variable, the type of which chosen through roulette
	// search on the cumulative weights of ActionTypeWeights,
	// resolved from ActionTypes to Actions using some strength
	// based on ActionStrengths.
	actions := make([]Action, 0)
	cum := CumulativeWeights(opt.ActionTypeWeights)
	for i := 0; i < len(e); i++ {
		for j := 0; j < opt.ActionCount.Poll(); j++ {
			k := algorithms.CumWeightedChooseOne(cum)
			a := opt.ActionTypes[k](e[i], opt.ActionStrengths[k].Poll())
			actions = append(actions, a)
		}
	}

	// scramble actions
	for i := 0; i < len(actions); i++ {
		j := rand.Intn(len(actions))
		actions[i], actions[j] = actions[j], actions[i]
	}

	splitIndex := int(math.Ceil(float64(len(actions)) * opt.PassiveRatio))

	gg.Actions = actions[0:splitIndex]
	gg.Passives = actions[splitIndex : len(actions)+1]

	gg.Init = make(map[int]float64)
	// Choose some number of variables to initialize
	// at game start
	l := opt.InitSize.Poll()
	for i := 0; i < l; i++ {
		j := rand.Intn(len(e))
		gg.Init[j] = opt.EnvVal.Poll()
	}

	gg.Goal = make(map[int]float64)

	// Simulate some actions on the environment . . .
	l := opt.GoalDistance.Poll()
	for i := 0; i < l; i++ {
		//Perform some action
		gg.Actions[rand.Intn(len(gg.Actions))]()
		//Perform all passive actions
		for j, p := range gg.Passives {
			p()
		}
	}

	// . . . and pull random elements from said environment
	// to determine the goal state.
	l := opt.GoalSize.Poll()
	for i := 0; i < l; i++ {
		j := rand.Intn(len(e))
		gg.Goal[j] = *e[j]
	}

	e.SetAll(0.0)
	gg.Environment = e

	return gg

}
