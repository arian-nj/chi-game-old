package bot

import (
	"context"
	"time"

	commonapp "github.com/arian-nj/chibazi/internals/common_app"

	"golang.org/x/exp/slog"
	"gopkg.in/telebot.v4"
)

type Application struct {
	*commonapp.CommonApp
}

func NewApplication(common *commonapp.CommonApp) *Application {
	return &Application{
		CommonApp: common,
	}
}

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
	b.Handle(telebot.OnInlineResult, app.inlineResultFeedbackHandler)
	b.Handle(telebot.OnCallback, app.callbackRouter)

	b.Handle(telebot.OnText, app.welcomeHandler)
	b.Handle("/start", app.welcomeHandler)
	b.Handle("/stat", app.statHandler)

	b.Handle(PlayWithFriendsButtonText, app.PlayWithFriendsHandler)
	b.Handle(PlayWithRandomPlayerText, app.PlayWithRandomPlayerHandler)

	b.Handle(Xo3x3ButtonText, app.PlayRandomXO3X3Handler)
	b.Handle(Xo5x5ButtonText, app.PlayRandomXO5X5Handler)

	go app.ClearDeadGamesCron()
	go app.MakeMatches()
	b.Start()
}

func (app *Application) ClearDeadGamesCron() {
	for {
		ExpireTime := 2 * time.Minute
		nowTime := time.Now()
		for key, hub := range app.Lobby.XOGames {
			if nowTime.Sub(hub.CreatedAt) > ExpireTime {
				delete(app.Lobby.XOGames, key)
			}
		}
		for key, hub := range app.Lobby.DotBox {
			if nowTime.Sub(hub.CreatedAt) > ExpireTime {
				delete(app.Lobby.DotBox, key)
			}
		}

		time.Sleep(1 * time.Minute)
	}
}

func (app *Application) addUserMiddleware(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		go func() {
			user := c.Sender()
			if user == nil {
				slog.Error("User is nil")
				return
			}
			err := app.Queries.CreateUser(context.Background(), int(user.ID))
			if err != nil {
				slog.Error("Failed to create user", "err", err)
				return
			}
		}()
		return next(c)
	}
}

const (
	PlayWithFriendsButtonText = "🎮 بازی تو پیوی یا گروه🫂"
	PlayWithRandomPlayerText  = "🎲 بازی با ناشناس 🕹"
)

const (
	Xo3x3ButtonText = "بازی دوز ۳ در ۳❌"
	Xo5x5ButtonText = "بازی دوز ۵ در ۵⭕️"
)

const (
	MainKeyboardButtonText = "صفحه اصلی🏠"
)
