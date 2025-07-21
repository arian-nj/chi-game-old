package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/arian-nj/chibazi/api"
	"github.com/arian-nj/chibazi/bot"
	"github.com/arian-nj/chibazi/db"
	"github.com/arian-nj/chibazi/internals/config"
	sharedapp "github.com/arian-nj/chibazi/internals/shared_app"
)

func main() {
	cfg, err := config.ParseConfig()
	if err != nil {
		panic(err)
	}

	err = db.Migrate(cfg.DatabseUrl)
	if err != nil {
		slog.Error("Failed to migrate database", "err", err)
		return
	}
	sharedApp := sharedapp.NewSharedApp(cfg)

	err = sharedApp.ConfigureDatabase()
	if err != nil {
		slog.Error("Failed to configure Database", "err", err)
		return
	}

	defer sharedApp.Conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	err = sharedApp.Conn.Ping(ctx)
	if err != nil {
		slog.Error("Failed to connect to Database", "err", err)
		return
	}
	cancel()

	slog.Info("Connected to Database")

	parentCtx, cancel := context.WithCancel(context.Background())

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	sharedApp.Wg.Add(1)
	go api.RunApi(sharedApp, parentCtx)

	sharedApp.Wg.Add(1)
	go bot.RunBot(sharedApp, parentCtx)

	<-quit
	cancel()
	sharedApp.Wg.Wait()
}
