package game

import (
	"context"

	sessionv1 "github.com/arian-nj/chibazi/backend/gen/session/v1"
	"github.com/arian-nj/chibazi/backend/internals/socket"
	"gopkg.in/telebot.v4"
)

type Game interface {
	AddPlayer(id int, name string, tgId int, socket *socket.Socket)

	SocketRouter(session *sessionv1.GameMessage, playerId int)
	CallBackRouter(c telebot.Context) error

	StartGame() error
	GetContext() context.Context

	GetGameData() *GameData

	SubToTelegram(ID int, bot *telebot.Bot, ViaMessageId string)
	SubToSocket(ID int, newSocket *socket.Socket) func()
}

type GameData struct {
	StartText string
	RulesText string
	MaxPlayer int
}

func NewGameData(startText, rulesText string, maxPlayer int) *GameData {
	return &GameData{
		StartText: startText,
		RulesText: rulesText,
		MaxPlayer: maxPlayer,
	}
}
