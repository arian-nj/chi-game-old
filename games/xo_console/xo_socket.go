package xoconsole

import (
	"encoding/json"
	"log/slog"

	"github.com/arian-nj/chibazi/internals/socket"
	"github.com/arian-nj/chibazi/internals/xo_core"
)

const (
	StartActionType socket.ActionType = "start"
	MoveActionType  socket.ActionType = "move"
)

func (g *XOGame) StartSocket() error {
	startAction := socket.NewGameAction(StartActionType, "")
	data_byte, err := json.Marshal(startAction)
	if err != nil {
		return err
	}
	for _, player := range g.Players {

		if player.Socket != nil {
			err := player.Socket.SendNewEvent(socket.GameEventType, string(data_byte))
			if err != nil {
				slog.Error("can't send new event to player", "err", err)
			}
		}
	}
	return nil
}

type MoveAction struct {
	MoveIndex int `json:"index"`
	CellType  int `json:"value"`
}

func (g *XOGame) BrodcastNewMove(moveIndex int, cellType xo_core.Cell) {
	for _, player := range g.Players {
		if player.Socket == nil {
			continue
		}

		err := player.Socket.SendNewAction(
			MoveActionType, MoveAction{
				MoveIndex: moveIndex,
				CellType:  int(cellType),
			})
		if err != nil {
			slog.Error("error broadcasting new move", "error", err)
		}
	}
}
