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
		endGameCommand = NewEndGameCommand(move.Game, move.Game.CurrentPlayer(), game.OpponentPlayer(), END_GAME_FULL)
	} else if !game.Board.IsAnyCellEmpty() {
		endGameCommand = NewEndGameCommand(move.Game, nil, nil, END_GAME_TIE)
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

type EndGameReason int

const (
	END_GAME_TIE EndGameReason = iota
	END_GAME_TIMEOUT
	END_GAME_FULL
)

type EndGameCommand struct {
	reason EndGameReason
	Winner *XoPlayer
	Loser  *XoPlayer
	Game   *XOState
}

func NewEndGameCommand(game *XOState, winner *XoPlayer, loser *XoPlayer, endGameReason EndGameReason) *EndGameCommand {
	return &EndGameCommand{
		reason: endGameReason,
		Winner: winner,
		Loser:  loser,
		Game:   game,
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
