package gamebot

import (
	commonapp "github.com/arian-nj/ultrun/internals/common_app"
	"github.com/arian-nj/ultrun/internals/config"
	"gopkg.in/telebot.v4"
)

type Application struct {
	*commonapp.CommonApp
	Lobby *Lobby
	Bot   *telebot.Bot
}

func NewApplication(conf *config.Config) *Application {
	return &Application{
		CommonApp: commonapp.NewCommon(conf),
		Lobby:     NewLobby(),
	}
}
