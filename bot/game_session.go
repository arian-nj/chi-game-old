package bot

import (
	consoleplayer "github.com/arian-nj/chibazi/internals/console_player"
	gametype "github.com/arian-nj/chibazi/internals/game_type"
)

type Game interface {
	Players() []*consoleplayer.ConsolePlayer
}

type GameSession struct {
	Gametype  gametype.GameType
	GameState Game
}
