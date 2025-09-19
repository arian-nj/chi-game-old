package conn4

import (
	"log/slog"

	conn4_core "github.com/arian-nj/chibazi/backend/games/conn4/core"
)

type StartCommand struct {
	Game *Conn4State
}

func NewStartCommand(game *Conn4State) *StartCommand {
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
	reason  EndGameReason
	Winner  *Conn4Player
	Loser   *Conn4Player
	Game    *Conn4State
	WinLine []int
}

func NewEndGameCommand(game *Conn4State, winner *Conn4Player, loser *Conn4Player, endGameReason EndGameReason, winLine []int) *EndGameCommand {
	return &EndGameCommand{
		reason:  endGameReason,
		Winner:  winner,
		Loser:   loser,
		Game:    game,
		WinLine: winLine,
	}
}

func (endGame *EndGameCommand) Execute() {
	endGame.Game.CancelGame()
}

type SyncTimeCommand struct {
	Game *Conn4State
}

func NewSyncTimeCommand(game *Conn4State) *SyncTimeCommand {
	return &SyncTimeCommand{
		Game: game,
	}
}

func (syncTime *SyncTimeCommand) Execute() {
}

type PlayCommand struct {
	PlayerID    int
	ColumnIndex int
	MoveType    conn4_core.Cell
	Game        *Conn4State
}

func NewPlayCommand(game *Conn4State, rowIndex int, moveType conn4_core.Cell, playerID int) *PlayCommand {
	return &PlayCommand{
		PlayerID:    playerID,
		ColumnIndex: rowIndex,
		MoveType:    moveType,
		Game:        game,
	}
}

func (move *PlayCommand) Execute() {
	currentPlayer := move.Game.CurrentPlayer()
	if currentPlayer.ID != move.PlayerID {
		slog.Error("play command recieved wrongly")
		return
	}
	game := move.Game

	idx, dropOk := game.Board.DropPiece(move.ColumnIndex, move.MoveType)
	if !dropOk {
		slog.Error("drop piece failed")
		return
	}

	hasWon, winLine := game.Board.HasWon(idx)

	var endGameCommand *EndGameCommand = nil
	if hasWon {
		endGameCommand = NewEndGameCommand(move.Game, currentPlayer, game.OpponentPlayer(), END_GAME_FULL, winLine)
	} else if !game.Board.IsAnyCellEmpty() {
		endGameCommand = NewEndGameCommand(move.Game, nil, nil, END_GAME_TIE, []int{})
	}

	if endGameCommand != nil {
		game.InjectCommand(endGameCommand)
		return
	}

	game.nextPlayer()
}
