package xo_core

type GameCommand interface {
	Execute()
}

type MoveCommand struct {
	PlayerID int64
	Pos      int
}

func (mv *MoveCommand) Execute() {

}
