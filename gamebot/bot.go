package gamebot

import (
	"fmt"
	"strings"
	"time"

	"github.com/arian-nj/ultrun/internals/config"
	"gopkg.in/telebot.v4"
)

func (app *Application) RunBot(cfg *config.Config) error {
	pref := telebot.Settings{
		Token:  cfg.BotToken,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}
	b, err := telebot.NewBot(pref)
	if err != nil {
		return fmt.Errorf("new error %w", err)
	}
	app.Bot = b

	b.Use(app.addUserMiddleware)
	b.Handle("/start", app.helloHandler)
	b.Handle(telebot.OnText, app.onMessage)
	b.Handle(telebot.OnQuery, app.onInlineQueryHandler)
	b.Handle(telebot.OnInlineResult, app.onInlineResult)
	b.Handle(telebot.OnCallback, app.onCallback)

	b.Start()
	return nil
}

var (
	selector = &telebot.ReplyMarkup{
		InlineKeyboard: [][]telebot.InlineButton{
			{
				{
					Text:                  "بازی با دوستان",
					InlineQueryChosenChat: &telebot.SwitchInlineQuery{AllowUserChats: true, AllowGroupChats: true},
				},
			},
		},
	}
)

func (app *Application) helloHandler(c telebot.Context) error {
	return c.Send(
		`خوش اومدید 👋
دکمه بازی با دوستان رو بزن تا تو هر چت یا گروهی با دوستات بازی کنی
	`, selector)
}

func (app *Application) onMessage(c telebot.Context) error {
	return c.Send(c.Message().Text)
}

func (app *Application) onCallback(c telebot.Context) error {
	callback := c.Callback()
	if callback.Data == "join_hub" {
		app.JoinHubHandler(c)
	} else if strings.HasPrefix(callback.Data, "xo_") {
		return app.XOCallbackHandlers(c)
	}
	return nil
}
