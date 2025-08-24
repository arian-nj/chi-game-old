package xo

type XoSubscriber interface {
	Update(state *XOGame, command Command)
}

type Command interface {
	Execute(game *XOGame)
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

func (mv *MoveCommand) Execute(game *XOGame) {
	game.Board.SetCell(mv.Pos, mv.MoveType)

	hasWon := game.Board.HasWon(mv.Pos)

	var endGameCommand *EndGameCommand = nil
	if hasWon {
		endGameCommand = NewEndGameCommand(game.GetCurrentPlayer(), "")
	} else if !game.Board.IsAnyCellEmpty() {
		endGameCommand = NewEndGameCommand(nil, "")
	}

	if endGameCommand != nil {
		game.InjectCommand(endGameCommand)
		return
	}
	game.NextPlayer()
	// return g.EditDuringGameBoard(c.Bot())
}

type StartCommand struct{}

func NewStartCommand() *StartCommand {
	return &StartCommand{}
}

func (mv *StartCommand) Execute(game *XOGame) {

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

func (mv *EndGameCommand) Execute(game *XOGame) {
	game.CancelGame()
}
