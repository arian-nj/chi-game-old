package gamebot

import (
	commonapp "github.com/arian-nj/ultrun/internals/common_app"
	"github.com/arian-nj/ultrun/internals/config"
)

type Application struct {
	*commonapp.CommonApp
}

func NewApplication(conf *config.Config) *Application {
	return &Application{
		CommonApp: commonapp.NewCommon(conf),
	}
}
