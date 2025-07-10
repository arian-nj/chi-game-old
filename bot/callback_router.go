package bot

import (
	"gopkg.in/telebot.v4"
	"strconv"
)

func (app *Application) callbackRouter(c telebot.Context) error {
	callback := c.Callback()
	messageId := c.Callback().MessageID
	if messageId == "" {
		messageId = strconv.Itoa(int(callback.Sender.ID))
	}

	gameSession, has_game := app.GameSessions[messageId]
	if has_game {
		return gameSession.GameState.CallbackHandler(c)
	}
	return c.RespondText("هیچی")
}
