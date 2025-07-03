package gamebot

import (
	"context"
	"time"

	commonapp "github.com/arian-nj/chibazi/internals/common_app"
	"golang.org/x/exp/slog"
	"gopkg.in/telebot.v4"
)

func RunBot(commonapp *commonapp.CommonApp, ctx context.Context) {
	defer commonapp.Wg.Done()
	app := NewApplication(commonapp)

	pref := telebot.Settings{
		Token:  commonapp.Config.BotToken,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}
	b, err := telebot.NewBot(pref)
	if err != nil {
		slog.Error("new error %w", err)
		return
	}
	go func() {
		<-ctx.Done()
		slog.Info("Shutting down Bot ...")
		b.Stop()
		slog.Info("Bot is shut down")
	}()
	app.Bot = b

	b.Use(app.addUserMiddleware)

	b.Handle(telebot.OnQuery, app.inlineQueryHandler)
	b.Handle(telebot.OnInlineResult, app.inlineResultHandler)
	b.Handle(telebot.OnCallback, app.callbackHandler)

	b.Handle(telebot.OnText, app.welcomeHandler)
	b.Handle("/start", app.welcomeHandler)
	b.Handle("/stat", app.statHandler)

	go app.ClearGamesCron()
	b.Start()
}
func (app *Application) ClearGamesCron() {
	for {
		nowTime := time.Now()
		for key, hub := range app.Lobby.Hubs {
			if nowTime.Sub(hub.CreatedAt) > 30*time.Minute {
				delete(app.Lobby.Hubs, key)
			}
		}
		time.Sleep(1 * time.Minute)
	}
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

func (app *Application) welcomeHandler(c telebot.Context) error {
	return c.Send(
		`خوش اومدید 👋
دکمه بازی با دوستان رو بزن تا تو هر چت یا گروهی با دوستات بازی کنی
	`, selector)
}
