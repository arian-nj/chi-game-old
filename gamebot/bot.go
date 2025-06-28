package gamebot

import (
	"fmt"
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

	b.Handle("/start", app.helloHandler)
	b.Handle(telebot.OnText, app.onMessage)
	b.Start()
	return nil
}

func (app *Application) helloHandler(c telebot.Context) error {
	return c.Send(`Welcone to Ultrun Bot`)
}

func (app *Application) onMessage(c telebot.Context) error {
	return c.Send(c.Message().Text)
}
