package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/arian-nj/chibazi/api"
	"github.com/arian-nj/chibazi/bot"
	"github.com/arian-nj/chibazi/database"
	"github.com/arian-nj/chibazi/db"
	gamesessions "github.com/arian-nj/chibazi/game_sessions"
	"github.com/arian-nj/chibazi/internals/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	var (
		Config      *config.Config
		Queries     *database.Queries
		Conn        *pgxpool.Pool
		Wg          *sync.WaitGroup          = &sync.WaitGroup{}
		AllSessions *gamesessions.AllSession = &gamesessions.AllSession{
			Sessions: map[string]*gamesessions.GameSession{},
			Mutex:    sync.Mutex{},
		}
	)

	Config, err := config.ParseConfig()
	if err != nil {
		panic(err)
	}

	err = db.Migrate(Config.DatabseUrl)
	if err != nil {
		slog.Error("Failed to migrate database", "err", err)
		return
	}

	Conn, err = pgxpool.New(context.Background(), Config.DatabseUrl)
	if err != nil {
		slog.Error("can not make a new connection ", "err", err)
		return
	}
	defer Conn.Close()

	Queries = database.New(Conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	err = Conn.Ping(ctx)
	if err != nil {
		slog.Error("Failed to connect to Database", "err", err)
		cancel()
		return
	}

	slog.Info("Connected to Database")

	parentCtx, cancel := context.WithCancel(context.Background())

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	app := api.NewApiApplication(Config, Queries, AllSessions)
	go app.RunApi(parentCtx, Wg)

	botApp := bot.NewBotApplication(Config, Queries, AllSessions)
	go botApp.RunBot(parentCtx, Wg)

	<-quit
	cancel()
	Wg.Wait()
}
