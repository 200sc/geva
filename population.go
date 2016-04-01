package neural

type Population []Network

type PopulationOptions struct {
	generationOptions *NetworkGenerationOptions
	size int
	selectionMethod *PopulationSelectionMethod
	crossoverMethod *PopulationCrossoverMethod
}

type PopulationSelectionMethod interface {
	Select(p_p *Population) *Population
}

type PopulationCrossoverMethod interface {
	Crossover(p_p *Population) *Population
}

type GreedySelection struct {

}

func (gs_p *GreedySelection) Select(p_p *Population) *Population {
	//gs := *gs_p
	//p := *p_p

	return p_p
}

type TournamentSelection struct {

}

func (ts_p *TournamentSelection) Select(p_p *Population) *Population {
	//ts := *ts_p
	//p := *p_p

	return p_p
}


type ProbabilisticSelection struct {

}

func (ps_p *ProbabilisticSelection) Select(p_p *Population) *Population {
	//ps := *ps_p
	//p := *p_p

	return p_p
}
