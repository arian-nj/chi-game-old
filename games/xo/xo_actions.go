package xo

type XoSubscriber interface {
	Update(state *XOGame, action Action)
}

type Action interface {
	Execute(game *XOGame)
}

type PlayAction struct {
	Pos      int
	MoveType Cell
}

func NewPlayAction(pos int, moveType Cell) *PlayAction {
	return &PlayAction{
		Pos:      pos,
		MoveType: moveType,
	}
}

func (mv *PlayAction) Execute(game *XOGame) {
	game.Board.SetCell(mv.Pos, mv.MoveType)

	hasWon := game.Board.HasWon(mv.Pos)

	var endGameAction *EndGameAction = nil
	if hasWon {
		endGameAction = NewEndGameAction(game.GetCurrentPlayer(), "")
	} else if !game.Board.IsAnyCellEmpty() {
		endGameAction = NewEndGameAction(nil, "")
	}

	if endGameAction != nil {
		game.InjectAction(endGameAction)
		return
	}
	game.NextPlayer()
	// return g.EditDuringGameBoard(c.Bot())
}

type StartAction struct{}

func NewStartAction() *StartAction {
	return &StartAction{}
}

func (mv *StartAction) Execute(game *XOGame) {

}

type EndGameAction struct {
	Winner *XoPlayer
	Text   string
}

func NewEndGameAction(winner *XoPlayer, text string) *EndGameAction {
	return &EndGameAction{
		Winner: winner,
		Text:   text,
	}
}

func (mv *EndGameAction) Execute(game *XOGame) {
	game.CancelGame()
}
