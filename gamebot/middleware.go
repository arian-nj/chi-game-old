package gamebot

import (
	"context"
	"log/slog"

	"gopkg.in/telebot.v4"
)

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
