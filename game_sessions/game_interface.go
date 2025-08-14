package gamesessions

import (
	"context"

	"github.com/arian-nj/chibazi/internals/socket"
	"gopkg.in/telebot.v4"
)

type Game interface {
	AddPlayer(name string, tgId int, socket *socket.Socket)
	CallbackHandler(c telebot.Context) error
	GetContext() context.Context
	StartGame() error
	SendJoinPanelAddSender(telebot.Context) error
	SetPlayerSocket(tgId int, socket *socket.Socket)
}
