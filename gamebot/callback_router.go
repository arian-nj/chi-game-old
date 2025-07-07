package gamebot

import (
	"strings"

	"github.com/arian-nj/chibazi/games/dotbox_console"
	xoconsole "github.com/arian-nj/chibazi/games/xo_console"
	gametype "github.com/arian-nj/chibazi/internals/game_type"
	"gopkg.in/telebot.v4"
)

func (app *Application) callbackRouter(c telebot.Context) error {
	callback := c.Callback()
	messageId := c.Callback().MessageID

	if after, has_prefix := strings.CutPrefix(callback.Data, string(gametype.XOGameType3X3)+"_"); has_prefix {
		xoGame, has_game := app.Lobby.XOGames[messageId]
		if has_game {
			return xoGame.XOCallbackHandlers(c, after)
		}
	} else if after, has_prefix := strings.CutPrefix(callback.Data, string(gametype.XOGameType5X5)+"_"); has_prefix {
		xoGame, has_game := app.Lobby.XOGames[messageId]
		if has_game {
			return xoGame.XOCallbackHandlers(c, after)
		}
	} else if after, has_prefix := strings.CutPrefix(callback.Data, string(gametype.DotBoxGameType)+"_"); has_prefix {
		dotBoxGame, has_game := app.Lobby.DotBox[messageId]
		if has_game {
			return dotBoxGame.CallbackHandlers(c, after)
		}
	}
	return c.RespondText("هیچی")
}

func (app *Application) inlineResultFeedbackHandler(c telebot.Context) error {
	resultId := c.InlineResult().ResultID
	messageID := c.InlineResult().MessageID

	switch gametype.GameType(resultId) {
	case gametype.XOGameType3X3:
		newXOGame := xoconsole.NewXOGame(gametype.XOGameType3X3, app.CommonApp)
		app.Lobby.XOGames[messageID] = newXOGame
		newXOGame.MessageId = messageID
		return newXOGame.SendJoinPanel(c)

	case gametype.XOGameType5X5:
		newXOGame := xoconsole.NewXOGame(gametype.XOGameType5X5, app.CommonApp)
		app.Lobby.XOGames[messageID] = newXOGame
		newXOGame.MessageId = messageID
		return newXOGame.SendJoinPanel(c)

	case gametype.DotBoxGameType:
		newDotbox := dotbox_console.NewDotBoxGame(app.CommonApp)
		app.Lobby.DotBox[messageID] = newDotbox
		newDotbox.MessageId = messageID
		return newDotbox.SendJoinPanel(c)
	}
	return c.RespondAlert("این بازیرو ندارم!")
}
