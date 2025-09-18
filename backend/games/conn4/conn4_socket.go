package conn4

import sessionv1 "github.com/arian-nj/chibazi/backend/gen/session/v1"

func (game *Conn4State) SocketRouter(newGameMsg *sessionv1.GameMessage, playerId int) {

	newXoMessage := newGameMsg.GetXo()
	switch newXoMessage.Payload.(type) {
	// case *xo_gamev1.XoGameMessage_Play:
	// playData := newXoMessage.GetPlay()
	// game.PlayHandlerSocket(playData, playerId)
	}
}
