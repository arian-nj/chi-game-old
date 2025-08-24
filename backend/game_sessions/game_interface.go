package gamesessions

import (
	"context"

	sessionv1 "github.com/arian-nj/chibazi/backend/gen/session/v1"
	"github.com/arian-nj/chibazi/backend/internals/socket"
	"gopkg.in/telebot.v4"
)

type Game interface {
	AddPlayer(id int, name string, tgId int, socket *socket.Socket)
	CallBackRouter(c telebot.Context) error
	SocketRouter(session *sessionv1.GameMessage, playerId int)
	GetContext() context.Context
	StartGame() error
	SendJoinPanelAddSender(telebot.Context) error
	SetPlayerSocket(tgId int, socket *socket.Socket)
}
