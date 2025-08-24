package api

import (
	"context"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/arian-nj/chibazi/backend/database"
	gamesessions "github.com/arian-nj/chibazi/backend/game_sessions"
	"github.com/arian-nj/chibazi/backend/internals/config"
	matchmaking "github.com/arian-nj/chibazi/backend/match_making"
)

type ApiApplication struct {
	Config      *config.Config
	Queries     *database.Queries
	AllSessions *gamesessions.AllSession
	MatchMaking *matchmaking.MatchMaking
}

func NewApiApplication(config *config.Config, queries *database.Queries, AllSession *gamesessions.AllSession, matchMaking *matchmaking.MatchMaking) *ApiApplication {
	return &ApiApplication{
		Config:      config,
		Queries:     queries,
		AllSessions: AllSession,
		MatchMaking: matchMaking,
	}
}

func (app *ApiApplication) RunApi(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	router := app.createRouter()

	srv := &http.Server{
		Addr:         ":8383",
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		slog.Info("Starting server on " + srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server listen error", "err", err)
			return
		}
	}()

	<-ctx.Done()

	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server forces to shutdown ", "err", err.Error())
	}
	slog.Info("Server exited")
}
