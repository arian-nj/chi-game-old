package xo

import "log/slog"

type XoSubscriber interface {
	Update(state *XOState, command Command)
}

type Command interface {
	Execute(game *XOState)
}

type MoveCommand struct {
	PlayerID int
	Pos      int
	MoveType Cell
}

func NewPlayCommand(pos int, moveType Cell, playerID int) *MoveCommand {
	return &MoveCommand{
		PlayerID: playerID,
		Pos:      pos,
		MoveType: moveType,
	}
}

func (move *MoveCommand) Execute(game *XOState) {
	if game.CurrentPlayer().ID != move.PlayerID {
		slog.Error("play commad recieved wrongly")
		return
	}
	game.Board.SetCell(move.Pos, move.MoveType)

	hasWon := game.Board.HasWon(move.Pos)

	var endGameCommand *EndGameCommand = nil
	if hasWon {
		endGameCommand = NewEndGameCommand(game.CurrentPlayer(), "")
	} else if !game.Board.IsAnyCellEmpty() {
		endGameCommand = NewEndGameCommand(nil, "")
	}

	if endGameCommand != nil {
		game.injectCommand(endGameCommand)
		return
	}
	game.nextPlayer()
	// return g.EditDuringGameBoard(c.Bot())
}

type StartCommand struct{}

func NewStartCommand() *StartCommand {
	return &StartCommand{}
}

func (start *StartCommand) Execute(game *XOState) {

}

type EndGameCommand struct {
	Winner *XoPlayer
	Text   string
}

func NewEndGameCommand(winner *XoPlayer, text string) *EndGameCommand {
	return &EndGameCommand{
		Winner: winner,
		Text:   text,
	}
}

func (endGame *EndGameCommand) Execute(game *XOState) {
	game.CancelGame()
}

type SyncTimeCommand struct {
}

func NewSyncTimeCommand() *SyncTimeCommand {
	return &SyncTimeCommand{}
}

func (syncTime *SyncTimeCommand) Execute(game *XOState) {
}
