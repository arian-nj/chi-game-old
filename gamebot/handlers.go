package gamebot

import (
	"context"
	"fmt"
	"strings"

	"github.com/arian-nj/ultrun/database"
	"gopkg.in/telebot.v4"
)

func (app *Application) inlineQueryHandler(c telebot.Context) error {
	results := make(telebot.Results, 1)
	xoResult := &telebot.ArticleResult{
		Title:       "دوز بازی",
		Description: "رو من کلیک کن",
		Text:        ticStartText,
	}
	xoResult.ParseMode = telebot.ModeMarkdownV2

	xoResult.ReplyMarkup = startInlineKeyboard

	xoResult.SetResultID(string(TicTacToe))
	results[0] = xoResult

	return c.Answer(&telebot.QueryResponse{
		Results:   results,
		CacheTime: 0,
	})
}

func (app *Application) inlineResultHandler(c telebot.Context) error {
	resultId := c.InlineResult().ResultID
	switch GameType(resultId) {
	case TicTacToe:
		return app.ticInlineResultReciever(c)
	}
	return c.RespondAlert("این بازیرو ندارم!")
}

func (app *Application) callbackHandler(c telebot.Context) error {
	callback := c.Callback()
	if callback.Data == "join_hub" {
		app.JoinHubHandler(c)
	} else if strings.HasPrefix(callback.Data, string(TicTacToe)+"_") {
		return app.TTTCallbackHandlers(c)
	}
	return nil
}
func (app *Application) JoinHubHandler(c telebot.Context) error {
	callback := c.Callback()
	messageId := callback.MessageID
	hub, isFound := app.Lobby.Hubs[messageId]
	if !isFound {
		return c.RespondAlert("همچین بازی ای وجود نداره!!!")
	}

	player := NewHumanPlayer(int(callback.Sender.ID), callback.Sender.FirstName)
	for _, joinedPlayer := range hub.Players {
		if joinedPlayer.TgID == player.TgID {
			return c.RespondText("تو بازی هستی")
		}
	}

	isOk := hub.Game.JoinGame(app, player)
	if !isOk {
		return c.RespondAlert("جا نیست")
	}
	err := c.RespondText("اضافه شدی")
	if err != nil {
		return err
	}

	_, err = app.Queries.CreateHub(context.Background(), database.CreateHubParams{GameType: string(hub.Game.GetGameType()), TgID: int(callback.Sender.ID)})
	return err

}

func (app *Application) statHandler(c telebot.Context) error {
	count, err := app.Queries.CountHubs(context.Background())
	if err != nil {
		return err
	}

	c.Send(fmt.Sprintf("تعداد بازی ها: %d", count))
	return nil
}
