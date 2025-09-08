package xo

import (
	"log/slog"
)

type MoveCommand struct {
	PlayerID int
	Pos      int
	MoveType Cell
	Game     *XOState
}

func NewPlayCommand(game *XOState, pos int, moveType Cell, playerID int) *MoveCommand {
	return &MoveCommand{
		PlayerID: playerID,
		Pos:      pos,
		MoveType: moveType,
		Game:     game,
	}
}

func (move *MoveCommand) Execute() {
	if move.Game.CurrentPlayer().ID != move.PlayerID {
		slog.Error("play commad recieved wrongly")
		return
	}
	game := move.Game
	game.Board.SetCell(move.Pos, move.MoveType)

	hasWon := game.Board.HasWon(move.Pos)

	var endGameCommand *EndGameCommand = nil
	if hasWon {
		endGameCommand = NewEndGameCommand(move.Game, move.Game.CurrentPlayer(), "")
	} else if !game.Board.IsAnyCellEmpty() {
		endGameCommand = NewEndGameCommand(move.Game, nil, "")
	}

	if endGameCommand != nil {
		game.InjectCommand(endGameCommand)
		return
	}
	game.nextPlayer()
	// return g.EditDuringGameBoard(c.Bot())
}

type StartCommand struct {
	Game *XOState
}

func NewStartCommand(game *XOState) *StartCommand {
	return &StartCommand{
		Game: game,
	}
}

func (start *StartCommand) Execute() {
}

type EndGameCommand struct {
	Winner *XoPlayer
	Text   string
	Game   *XOState
}

func NewEndGameCommand(game *XOState, winner *XoPlayer, text string) *EndGameCommand {
	return &EndGameCommand{
		Winner: winner,
		Text:   text,
	}
}

func (endGame *EndGameCommand) Execute() {
	endGame.Game.CancelGame()
}

type SyncTimeCommand struct {
	Game *XOState
}

func NewSyncTimeCommand(game *XOState) *SyncTimeCommand {
	return &SyncTimeCommand{
		Game: game,
	}
}

func (syncTime *SyncTimeCommand) Execute() {
}
