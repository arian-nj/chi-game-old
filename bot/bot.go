package bot

import (
	"context"
	"sync"
	"time"

	commonapp "github.com/arian-nj/chibazi/internals/common_app"
	gametype "github.com/arian-nj/chibazi/internals/game_type"

	"golang.org/x/exp/slog"
	"gopkg.in/telebot.v4"
)

type Application struct {
	*commonapp.CommonApp

	GameSessions *AllSession
	MatchMaking  *MatchMaking
}

func NewApplication(common *commonapp.CommonApp) *Application {
	return &Application{
		CommonApp: common,
		GameSessions: &AllSession{
			Sessions: map[string]*GameSession{},
			Mutex:    sync.Mutex{},
		},
		MatchMaking: &MatchMaking{
			WaitingPlayers: map[gametype.GameType][]*Ticket{},
			Mutex:          sync.Mutex{},
		},
	}
}

func (app *Application) AddGameSession(key string, gs *GameSession) {
	app.GameSessions.Mutex.Lock()
	defer app.GameSessions.Mutex.Unlock()

	app.GameSessions.Sessions[key] = gs
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

	b.Handle("/start", app.welcomeHandler)
	b.Handle("/stat", app.statHandler)

	b.Handle(PlayWithFriendsButtonText, app.PlayWithFriendsHandler)
	b.Handle(PlayWithRandomPlayerText, app.PlayWithRandomPlayerHandler)

	b.Handle(Xo3x3ButtonText, app.PlayRandomXO3X3Handler)
	b.Handle(Xo5x5ButtonText, app.PlayRandomXO5X5Handler)

	b.Handle(CancelGameButtonText, app.CancelSearchingForGame)

	b.Handle(StopChatButtonText, app.StopChatHandler)

	b.Handle(telebot.OnText, app.textHandler)

	go app.ClearDeadGamesCron()
	go app.MakeMatches()
	b.Start()
}

func (app *Application) ClearDeadGamesCron() {
	for {
		nowTime := time.Now()
		for key, gameSession := range app.GameSessions.Sessions {
			if nowTime.Sub(gameSession.CreatedAt) > gameSession.ExpireDuaration {
				app.GameSessions.Mutex.Lock()
				delete(app.GameSessions.Sessions, key)
				app.GameSessions.Mutex.Unlock()
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

	CancelGameButtonText = "🔙 نگرد پشیمون شدم"
)

const (
	MainKeyboardButtonText = "صفحه اصلی🏠"
)

const (
	StopChatButtonText = "قطع چت ✂️"
)
