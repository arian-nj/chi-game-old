package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/arian-nj/chibazi/db"
	"github.com/arian-nj/chibazi/gamebot"
	commonapp "github.com/arian-nj/chibazi/internals/common_app"
	"github.com/arian-nj/chibazi/internals/config"
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
	commonApp := commonapp.NewCommon(cfg)

	err = commonApp.ConfigureDatabase()
	if err != nil {
		slog.Error("Failed to configure Database", "err", err)
		return
	}

	defer commonApp.Conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	err = commonApp.Conn.Ping(ctx)
	if err != nil {
		slog.Error("Failed to connect to Database", "err", err)
		return
	}
	cancel()

	slog.Info("Connected to Database")

	parentCtx, cancel := context.WithCancel(context.Background())

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// commonApp.Wg.Add(1)
	// go api.RunApi(commonApp, parentCtx)

	commonApp.Wg.Add(1)
	go gamebot.RunBot(commonApp, parentCtx)

	<-quit
	cancel()
	commonApp.Wg.Wait()
}
