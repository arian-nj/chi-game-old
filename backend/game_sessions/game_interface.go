package gamesessions

import (
	"context"

	sessionv1 "github.com/arian-nj/chibazi/backend/gen/session/v1"
	"github.com/arian-nj/chibazi/backend/internals/socket"
	"gopkg.in/telebot.v4"
)

type Game interface {
	AddPlayer(id int, name string, tgId int, socket *socket.Socket)
	SetPlayerSocket(ID int, socket *socket.Socket)
	SocketRouter(session *sessionv1.GameMessage, playerId int)

	CallBackRouter(c telebot.Context) error
	SendJoinPanelAddSender(telebot.Context) error

	StartGame() error
	GetContext() context.Context
}
