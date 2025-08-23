package xo

import (
	"encoding/json"
	"log/slog"

	"github.com/arian-nj/chibazi/internals/socket"
)

// FIXME: make it map handler with auto Converting inputs
func (game *XOGame) SocketRouter(actionRawData json.RawMessage, playerId int) {
	newGameAction := &socket.GameAction{}
	err := json.Unmarshal(actionRawData, newGameAction)
	if err != nil {
		slog.Error("can't unmarshal action raw data", "err", err)
		return
	}
	AData, ok := newGameAction.AData.([]byte)
	if !ok {
		slog.Error("can't use adata as byte", "data", AData)
		return
	}

	switch newGameAction.Action {
	case PlayActionType:
		newPlayInput := &PlayActionInput{}
		err := json.Unmarshal(AData, newPlayInput)
		if err != nil {
			slog.Error("can't unmarshal play action input", "error", err)
			return
		}
		game.PlayHandlerSocket(newPlayInput, playerId)
	}
}

type PlayActionInput struct {
	MoveIndex int `json:"index"`
}

func (game *XOGame) PlayHandlerSocket(playInput *PlayActionInput, playerID int) {
	if game.IsPlayersTurn(playerID) == false {
		slog.Error("not players turn")
		return
	}

	moveType := game.GetCurrentPlayer().Move

	isValid, errMsg := game.Board.IsMoveValid(playInput.MoveIndex, moveType)
	if !isValid {
		slog.Error("move is not valid", "err", errMsg)
		return
	}

	player := game.Find(playerID)
	if player == nil {
		slog.Error("can't find player in socjet play handler")
		return
	}
	playCommand := NewPlayCommand(playInput.MoveIndex, moveType, player.ID)
	game.PushCommand(playCommand)
}

type SocketListener struct{}

func (sl *SocketListener) Update(game *XOGame, command Command) {
	switch a := command.(type) {
	case *PlayCommand:
		sl.SocketBrodcastNewMove(game, a.Pos, a.MoveType, a.PlayerID)
	case *StartCommand:
	case *EndGameCommand:
		if a.Winner == nil {
		} else {
		}
	}
}

const (
	StartActionType socket.ActionType = "start"
	MoveActionType  socket.ActionType = "move"
	EndActionType   socket.ActionType = "end"

	PlayActionType socket.ActionType = "play"
)

func (g *XOGame) StartSocket() error {
	startAction := socket.NewGameAction(StartActionType, json.RawMessage{})
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

func (sl *SocketListener) SocketBrodcastNewMove(game *XOGame, moveIndex int, cellType Cell, playerID int) {
	for _, player := range game.Players {
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
