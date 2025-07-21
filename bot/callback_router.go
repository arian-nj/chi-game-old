package bot

import (
	"gopkg.in/telebot.v4"
	"strconv"
)

func (app *BotApplication) callbackRouter(c telebot.Context) error {
	callback := c.Callback()
	messageId := c.Callback().MessageID
	if messageId == "" {
		messageId = strconv.Itoa(int(callback.Sender.ID))
	}

	app.AllSessions.Mutex.Lock()
	gameSession, has_game := app.AllSessions.Sessions[messageId]
	app.AllSessions.Mutex.Unlock()

	if has_game {
		return gameSession.GameState.CallbackHandler(c)
	}
	return c.RespondText("هیچی")
}
