package api

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	sharedapp "github.com/arian-nj/chibazi/internals/shared_app"
)

type Application struct {
	*sharedapp.SharedApp
}

func NewApplication(sharedApp *sharedapp.SharedApp) *Application {
	return &Application{
		SharedApp: sharedApp,
	}
}

func RunApi(sharedApp *sharedapp.SharedApp, ctx context.Context) {
	defer sharedApp.Wg.Done()
	app := NewApplication(sharedApp)
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
