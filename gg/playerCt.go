package gg

type PlayerCount struct {
	devIndex    int
	playerCount int
}

type PlayerCounts []PlayerCount

func (a PlayerCounts) Len() int           { return len(a) }
func (a PlayerCounts) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a PlayerCounts) Less(i, j int) bool { return a[i].playerCount < a[j].playerCount }
