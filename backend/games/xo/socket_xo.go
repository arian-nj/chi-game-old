package xo

import (
	"log/slog"

	sessionv1 "github.com/arian-nj/chibazi/backend/gen/session/v1"
	xo_gamev1 "github.com/arian-nj/chibazi/backend/gen/xo_game/v1"
)

// FIXME: make it map handler with auto Converting inputs
func (game *XOState) SocketRouter(newGameMsg *sessionv1.GameMessage, playerId int) {
	newXoMessage := newGameMsg.GetXo()
	switch newXoMessage.Payload.(type) {
	case *xo_gamev1.XoGameMessage_Play:
		playData := newXoMessage.GetPlay()
		game.PlayHandlerSocket(playData, playerId)
	}
}

type SocketListener struct{}

func (sl *SocketListener) Update(game *XOState, command Command) {
	switch a := command.(type) {
	case *MoveCommand:
		sl.SocketBrodcastNewMove(game, a.Pos, a.MoveType, a.PlayerID)
	case *StartCommand:
	case *EndGameCommand:
		if a.Winner == nil {
		} else {
		}
	case *SyncTimeCommand:
		sl.SocketBrodcastSyncTime(game)
	}
}

func sendInvalidResponse(player *XoPlayer, errMsg string, cellIndex int32, cellType Cell) error {
	newSessionMsg := sessionv1.SessionMessage{
		Content: &sessionv1.SessionMessage_Game{
			Game: &sessionv1.GameMessage{
				Game: &sessionv1.GameMessage_Xo{
					Xo: &xo_gamev1.XoGameMessage{
						Payload: &xo_gamev1.XoGameMessage_PlayResponse{
							PlayResponse: &xo_gamev1.PlayResponse{
								IsValid: false,
								Reason:  errMsg,
								Move: &xo_gamev1.Move{
									CellIndex: cellIndex,
									CellValue: int32(cellType),
								},
							},
						},
					},
				},
			},
		},
	}

	return player.Socket.SendMessage(&newSessionMsg)
}

func (game *XOState) PlayHandlerSocket(playInput *xo_gamev1.Play, playerID int) {
	player := game.findByID(playerID)
	if player == nil {
		slog.Error("can't find player in socjet play handler")
		return
	}

	if game.CurrentPlayer().ID != playerID {
		sendInvalidResponse(player, "نوبت تو نیست", playInput.CellIndex, 0)
		return
	}

	moveType := game.CurrentPlayer().Move

	isValid, errMsg := game.Board.IsMoveValid(int(playInput.CellIndex), moveType)
	if !isValid {
		err := sendInvalidResponse(player, errMsg, playInput.CellIndex, 0)
		if err != nil {
			slog.Error("can't send invalid response")
		}
		return
	}

	playCommand := NewPlayCommand(int(playInput.CellIndex), moveType, player.ID)
	game.pushCommand(playCommand)
}

func (sl *SocketListener) SocketBrodcastNewMove(game *XOState, moveIndex int, cellType Cell, playerID int) {
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
										Move: &xo_gamev1.Move{
											CellIndex: int32(moveIndex),
											CellValue: int32(cellType),
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
										CellIndex: int32(moveIndex),
										CellValue: int32(cellType),
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
func (sl *SocketListener) SocketBrodcastSyncTime(gameState *XOState) {
	for _, player := range gameState.Players {
		newSessionMessage := sessionv1.SessionMessage{
			Content: &sessionv1.SessionMessage_Game{
				Game: &sessionv1.GameMessage{
					Game: &sessionv1.GameMessage_Xo{
						Xo: &xo_gamev1.XoGameMessage{
							Payload: &xo_gamev1.XoGameMessage_SyncTime{
								SyncTime: &xo_gamev1.Time{
									SpentTime: int32(player.SpentTime),
									TotalTime: int32(MaxAllowedTimeSecond),
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
