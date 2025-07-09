package bot

import (
	"strconv"
	"strings"

	gametype "github.com/arian-nj/chibazi/internals/game_type"
	"gopkg.in/telebot.v4"
)

func (app *Application) callbackRouter(c telebot.Context) error {
	callback := c.Callback()
	messageId := c.Callback().MessageID
	if messageId == "" {
		messageId = strconv.Itoa(int(callback.Sender.ID))
	}

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
