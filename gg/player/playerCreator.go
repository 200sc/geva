package player

type Creator interface {
	NewPlayer() Player
}
