package main

import (
	"context"
	"log/slog"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/arian-nj/chibazi/backend/api"
	"github.com/arian-nj/chibazi/backend/bot"
	"github.com/arian-nj/chibazi/backend/database"
	"github.com/arian-nj/chibazi/backend/db"
	gamesessions "github.com/arian-nj/chibazi/backend/game_sessions"
	"github.com/arian-nj/chibazi/backend/internals/config"
	matchmaking "github.com/arian-nj/chibazi/backend/match_making"
	"github.com/jackc/pgx/v5/pgxpool"
	"gopkg.in/telebot.v4"
)

// NOTE: use matchmaking as a package not to run code u fucking idiot

type GlobalVars struct {
	Config      *config.Config
	Queries     *database.Queries
	Conn        *pgxpool.Pool
	Wg          *sync.WaitGroup
	AllSessions *gamesessions.AllSession
	MatchMaking *matchmaking.MatchMaking
	Bot         *telebot.Bot
}

func NewGlobalVars() *GlobalVars {
	return &GlobalVars{
		AllSessions: &gamesessions.AllSession{
			Sessions: map[string]*gamesessions.GameSession{},
			Mutex:    sync.Mutex{},
		},
		Wg: &sync.WaitGroup{},
	}
}

func main() {
	// go func() {
	// 	log.Println(http.ListenAndServe("localhost:6060", nil))
	// }()
	GlobalVars := NewGlobalVars()

	var err error
	GlobalVars.Config, err = config.ParseConfig()
	if err != nil {
		panic(err)
	}

	err = db.Migrate(GlobalVars.Config.DatabseUrl)
	if err != nil {
		slog.Error("Failed to migrate database", "err", err)
		return
	}

	GlobalVars.Conn, err = pgxpool.New(context.Background(), GlobalVars.Config.DatabseUrl)
	if err != nil {
		slog.Error("can not make a new connection ", "err", err)
		return
	}
	defer GlobalVars.Conn.Close()

	GlobalVars.Queries = database.New(GlobalVars.Conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	err = GlobalVars.Conn.Ping(ctx)
	if err != nil {
		slog.Error("Failed to connect to Database", "err", err)
		cancel()
		return
	}
	GlobalVars.MatchMaking = matchmaking.NewMatchMaking(GlobalVars.AllSessions, GlobalVars.Queries)

	slog.Info("Connected to Database")

	parentCtx, cancel := context.WithCancel(context.Background())

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	app := api.NewApiApplication(GlobalVars.Config, GlobalVars.Queries, GlobalVars.AllSessions, GlobalVars.MatchMaking)
	go app.RunApi(parentCtx, GlobalVars.Wg)

	botApp := bot.NewBotApplication(GlobalVars.Config, GlobalVars.Queries, GlobalVars.AllSessions, GlobalVars.MatchMaking)
	bot, err := botApp.MakeBot()
	if err != nil {
		slog.Error("Failed to make bot", "err", err)
		return
	}
	GlobalVars.Bot = bot
	go botApp.RunBot(bot, parentCtx, GlobalVars.Wg)

	go ClearDeadGamesCron(GlobalVars.AllSessions)

	go GlobalVars.MakeMatches()

	<-quit
	cancel()
	GlobalVars.Wg.Wait()
}

func ClearDeadGamesCron(allSessions *gamesessions.AllSession) {
	for {
		nowTime := time.Now()
		for key, gameSession := range allSessions.Sessions {
			if nowTime.Sub(gameSession.CreatedAt) > gameSession.ExpireDuaration {
				allSessions.Mutex.Lock()
				delete(allSessions.Sessions, key)
				allSessions.Mutex.Unlock()
			}
		}
		time.Sleep(1 * time.Minute)
	}
}
