package gamebot

import (
	"strconv"

	"gopkg.in/telebot.v4"
)

func (app *Application) onInlineQueryHandler(c telebot.Context) error {
	results := make(telebot.Results, 1)
	xoResult := &telebot.ArticleResult{
		Title:       "دوز بازی",
		Description: "رو من کلیک کن",
		Text:        ticStartText,
	}
	xoResult.ParseMode = telebot.ModeMarkdownV2

	xoResult.ReplyMarkup = startInlineKeyboard

	xoResult.SetResultID(strconv.Itoa(1))
	results[0] = xoResult

	return c.Answer(&telebot.QueryResponse{
		Results:   results,
		CacheTime: 0,
	})
}

func (app *Application) onInlineResult(c telebot.Context) error {
	return app.ticInlineResultReciever(c)
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
	return c.RespondText("اضافه شدی")
}
