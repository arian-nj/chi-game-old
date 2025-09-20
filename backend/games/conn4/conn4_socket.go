package conn4

import (
	sessionv1 "github.com/arian-nj/chibazi/backend/gen/session/v1"
	"github.com/arian-nj/chibazi/backend/internals/commander"
)

func (game *Conn4State) SocketRouter(newGameMsg *sessionv1.GameMessage, playerId int) {

	newXoMessage := newGameMsg.GetXo()
	switch newXoMessage.Payload.(type) {
	// case *xo_gamev1.XoGameMessage_Play:
	// playData := newXoMessage.GetPlay()
	// game.PlayHandlerSocket(playData, playerId)
	}
}

type SocketListener struct {
	Player *Conn4Player
}

func (sl *SocketListener) Update(command commander.Command) {
	// switch c := command.(type) {
	// case *MoveCommand
	// sl.SocketBrodcastNewMove(c.Game, c.Pos, c.MoveType, c.PlayerID)
	// case *StartCommand:
	// case *EndGameCommand:
	// sl.SendEndGameSocket(c)
	// case *SyncTimeCommand:
	// sl.SocketBrodcastSyncTime(c.Game)
	// }
}

//
// func (sl *SocketListener) SocketBrodcastSyncTime(gameState *Conn4State) {
// 	allTimeMessages := []*sessionv1.SessionMessage{}
// 	for _, player := range gameState.Players {
// 		newSessionMessage := sessionv1.SessionMessage{
// 			Content: &sessionv1.SessionMessage_Game{
// 				Game: &sessionv1.GameMessage{
// 					Game: &conn4_gamev1.Conn4GameMessage,
// 				},
// 			},
// 		}
// 		allTimeMessages = append(allTimeMessages, &newSessionMessage)
// 	}
//
// 	player := gameState.findByID(sl.Player.ID)
// 	for _, timeEvent := range allTimeMessages {
// 		err := player.Socket.SendMessage(timeEvent)
// 		if err != nil {
// 			slog.Error("can't send new move", "err", err)
// 		}
// 	}
// }
