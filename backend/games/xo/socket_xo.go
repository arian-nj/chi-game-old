package xo

import (
	"encoding/json"
	"log/slog"

	sessionv1 "github.com/arian-nj/chibazi/backend/gen/session/v1"
	xo_gamev1 "github.com/arian-nj/chibazi/backend/gen/xo_game/v1"
	"github.com/arian-nj/chibazi/backend/internals/socket"
)

// FIXME: make it map handler with auto Converting inputs
func (game *XOGame) SocketRouter(newGameMsg *sessionv1.GameMessage, playerId int) {
	newXoMessage := newGameMsg.GetXo()
	switch newXoMessage.Payload.(type) {
	case *xo_gamev1.XoGameMessage_Play:
		playData := newXoMessage.GetPlay()
		game.PlayHandlerSocket(playData, playerId)
	}
}

type SocketListener struct{}

func (sl *SocketListener) Update(game *XOGame, command Command) {
	switch a := command.(type) {
	case *MoveCommand:
		sl.SocketBrodcastNewMove(game, a.Pos, a.MoveType, a.PlayerID)
	case *StartCommand:
	case *EndGameCommand:
		if a.Winner == nil {
		} else {
		}
	}
}

func (game *XOGame) PlayHandlerSocket(playInput *xo_gamev1.Play, playerID int) {
	if game.IsPlayersTurnID(playerID) == false {
		slog.Error("not players turn")
		return
	}

	moveType := game.GetCurrentPlayer().Move

	isValid, errMsg := game.Board.IsMoveValid(int(playInput.CellIndex), moveType)
	if !isValid {
		slog.Error("move is not valid", "err", errMsg)
		return
	}

	player := game.FindByID(playerID)
	if player == nil {
		slog.Error("can't find player in socjet play handler")
		return
	}
	playCommand := NewPlayCommand(int(playInput.CellIndex), moveType, player.ID)
	game.PushCommand(playCommand)
}

const (
	StartActionType socket.ActionType = "start"
	MoveActionType  socket.ActionType = "move"
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

		if player.ID == playerID {
			newSessionMessage := sessionv1.SessionMessage{
				Content: &sessionv1.SessionMessage_Game{
					Game: &sessionv1.GameMessage{
						Game: &sessionv1.GameMessage_Xo{
							Xo: &xo_gamev1.XoGameMessage{
								Payload: &xo_gamev1.XoGameMessage_PlayResponse{
									PlayResponse: &xo_gamev1.PlayResponse{
										IsValid: true,
										Play: &xo_gamev1.Play{
											CellIndex: int32(moveIndex),
										},
									},
								},
							},
						},
					},
				},
			}

			err := player.Socket.SendMessage(&newSessionMessage)
			if err != nil {
				slog.Error("can't send play response", "err", err)
			}
		} else {
			newSessionMessage := sessionv1.SessionMessage{
				Content: &sessionv1.SessionMessage_Game{
					Game: &sessionv1.GameMessage{
						Game: &sessionv1.GameMessage_Xo{
							Xo: &xo_gamev1.XoGameMessage{
								Payload: &xo_gamev1.XoGameMessage_Move{
									Move: &xo_gamev1.Move{
										PlayerId:  int32(playerID),
										CellIndex: int32(moveIndex),
									},
								},
							},
						},
					},
				},
			}
			err := player.Socket.SendMessage(&newSessionMessage)
			if err != nil {
				slog.Error("can't send new move", "err", err)
			}
		}
	}
}
