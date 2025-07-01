package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/arian-nj/ultrun/db"
	"github.com/arian-nj/ultrun/gamebot"
	"github.com/arian-nj/ultrun/internals/config"
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
	manApp := gamebot.NewApplication(cfg)

	err = manApp.ConfigureDatabase()
	if err != nil {
		slog.Error("Failed to configure Database", "err", err)
		return
	}

	defer manApp.Conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	err = manApp.Conn.Ping(ctx)
	if err != nil {
		slog.Error("Failed to connect to Database", "err", err)
		return
	}
	cancel()

	slog.Info("Connected to Database")

	err = manApp.RunBot(cfg)
	if err != nil {
		slog.Error("Failed to run bot", "err", err)
		return
	}

}
