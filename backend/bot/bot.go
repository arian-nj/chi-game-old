package bot

import (
	"context"
	"fmt"
	"log"
	"runtime/debug"
	"strconv"
	"sync"
	"time"

	"log/slog"

	"github.com/arian-nj/chibazi/backend/database"
	gamesessions "github.com/arian-nj/chibazi/backend/game_sessions"
	"github.com/arian-nj/chibazi/backend/internals/config"
	"github.com/arian-nj/chibazi/backend/internals/keybul"
	"github.com/arian-nj/chibazi/backend/internals/utils"
	matchmaking "github.com/arian-nj/chibazi/backend/match_making"
	"gopkg.in/telebot.v4"
)

type BotApplication struct {
	Config      *config.Config
	Queries     *database.Queries
	Bot         *telebot.Bot
	AllSessions *gamesessions.AllSession
	MatchMaking *matchmaking.MatchMaking
}

func NewBotApplication(conf *config.Config, queries *database.Queries, AllSession *gamesessions.AllSession, mamatchmaking *matchmaking.MatchMaking) *BotApplication {
	return &BotApplication{
		Config:      conf,
		Queries:     queries,
		MatchMaking: mamatchmaking,
		AllSessions: AllSession,
	}
}

func (app *BotApplication) RunBot(bot *telebot.Bot, ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	go func() {
		<-ctx.Done()
		slog.Info("Shutting down Bot ...")
		bot.Stop()
		slog.Info("Bot is shut down")
	}()

	bot.Start()
}

func (app *BotApplication) MakeBot() (*telebot.Bot, error) {
	pref := telebot.Settings{
		Token:  app.Config.BotToken,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}
	b, err := telebot.NewBot(pref)
	if err != nil {
		return nil, err
	}

	app.Bot = b

	b.Use(panicRecover)
	b.Use(app.addUserMiddleware)

	b.Handle(telebot.OnQuery, app.inlineQueryHandler)
	b.Handle(telebot.OnInlineResult, app.inlineResultFeedbackHandler)
	b.Handle(telebot.OnCallback, app.handleCallback)

	b.Handle("/start", app.welcomeHandler)
	b.Handle("/panel", app.statHandler)
	b.Handle("/me", app.meHandler)

	b.Handle(keybul.PlayWithFriendsButtonText, app.PlayWithFriendsHandler)
	b.Handle(keybul.PlayWithRandomPlayerText, app.PlayWithRandomPlayerHandler)

	b.Handle(keybul.Xo3x3ButtonText, app.PlayRandomXO3X3Handler)
	// b.Handle(keybul.Xo5x5ButtonText, app.PlayRandomXO5X5Handler)
	b.Handle(keybul.Conn4ButtonText, app.PlayRandomConn4Handler)

	b.Handle(keybul.CancelGameButtonText, app.CancelSearchingForGame)

	b.Handle(keybul.StopChatButtonText, app.StopChatHandler)

	b.Handle(telebot.OnText, app.textHandler)

	return b, nil
}

func (app *BotApplication) addUserMiddleware(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		utils.RunBackgroundTask(func() {
			user := c.Sender()
			if user == nil {
				slog.Error("User is nil")
				return
			}
			_, err := app.Queries.CreateTgUser(context.Background(), database.CreateTgUserParams{
				TgID: int(c.Sender().ID),
				Name: c.Sender().FirstName,
			})
			if err != nil {
				slog.Error("Failed to create user", "err", err)
				return
			}
		})
		return next(c)
	}
}
func panicRecover(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(c telebot.Context) (err error) {
		defer func() {
			if r := recover(); r != nil {
				// log.Printf("Recovered from panic: %v\n", r)
				log.Printf("panic: %v\n%s", r, debug.Stack())

				c.Send("An internal error occurred. Please try again later.")

				// Set err so it appears as if the handler returned an error
				err = fmt.Errorf("handler panic: %v", r)
			}
		}()

		return next(c)
	}
}

var IDtoTgIDMap = make(map[int64]int)

func GetUserID(tgId int64, queries *database.Queries) int {
	id, cacheSuccess := IDtoTgIDMap[tgId]
	if cacheSuccess {
		return id
	}
	user, err := queries.GetTgUserByTgID(context.Background(), int(tgId))
	if err != nil {
		slog.Info("can't find user", "tg_id", tgId, "err", err)
		return 0
	}
	IDtoTgIDMap[tgId] = user.ID
	return user.ID
}

func (app *BotApplication) handleCallback(c telebot.Context) error {
	callback := c.Callback()
	messageId := callback.MessageID
	if messageId == "" {
		messageId = strconv.Itoa(GetUserID(callback.Sender.ID, app.Queries))
	}

	gameSession, hasSession := app.AllSessions.Get(messageId)

	if hasSession {
		return gameSession.HandleCallback(c, app.Queries)
	}
	return c.RespondText("no active game")
}
