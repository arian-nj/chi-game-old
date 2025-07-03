package gamebot

import (
	commonapp "github.com/arian-nj/chibazi/internals/common_app"
)

type Application struct {
	*commonapp.CommonApp
	Lobby *Lobby
}

func NewApplication(common *commonapp.CommonApp) *Application {
	return &Application{
		CommonApp: common,
		Lobby:     NewLobby(),
	}
}
