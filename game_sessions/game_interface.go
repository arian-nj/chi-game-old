package gamesessions

import (
	"context"
	"encoding/json"

	"github.com/arian-nj/chibazi/internals/socket"
	"gopkg.in/telebot.v4"
)

type Game interface {
	AddPlayer(id int, name string, tgId int, socket *socket.Socket)
	CallBackRouter(c telebot.Context) error
	SocketRouter(gameActionData json.RawMessage, playerId int)
	GetContext() context.Context
	StartGame() error
	SendJoinPanelAddSender(telebot.Context) error
	SetPlayerSocket(tgId int, socket *socket.Socket)
}
