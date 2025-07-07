package gamebot

import (
	"context"
	"time"

	"github.com/arian-nj/chibazi/games/dotbox_console"
	xoconsole "github.com/arian-nj/chibazi/games/xo_console"
	commonapp "github.com/arian-nj/chibazi/internals/common_app"
	"github.com/arian-nj/chibazi/linedot"
	"golang.org/x/exp/slog"
	"gopkg.in/telebot.v4"
)

type Lobby struct {
	XOGames map[string]*xoconsole.XOGame
	DotBox  map[string]*dotbox_console.DotBoxGame
}

type Application struct {
	*commonapp.CommonApp
	Lobby        *Lobby
	DotLineGames map[string]*linedot.DotLineGame
}

func NewApplication(common *commonapp.CommonApp) *Application {
	return &Application{
		CommonApp: common,
		Lobby: &Lobby{
			XOGames: map[string]*xoconsole.XOGame{},
			DotBox:  map[string]*dotbox_console.DotBoxGame{},
		},
		// DotLineGames: map[string]*linedot.DotLineGame{},
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

	go app.ClearDeadGamesCron()
	b.Start()
}

func (app *Application) ClearDeadGamesCron() {
	for {
		ExpireTime := 15 * time.Minute
		nowTime := time.Now()
		for key, hub := range app.Lobby.XOGames {
			if nowTime.Sub(hub.CreatedAt) > ExpireTime {
				delete(app.Lobby.XOGames, key)
			}
		}
		for key, hub := range app.Lobby.DotBox {
			if nowTime.Sub(hub.CreatedAt) > ExpireTime {
				delete(app.Lobby.XOGames, key)
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
