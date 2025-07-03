package gamebot

import (
	"context"
	"fmt"
	"strings"

	"github.com/arian-nj/chibazi/database"
	"gopkg.in/telebot.v4"
)

func (app *Application) inlineQueryHandler(c telebot.Context) error {
	results := telebot.Results{}
	// XO result
	xoResult := &telebot.ArticleResult{
		Title:       "دوز بازی",
		Description: "رو من کلیک کن",
		Text:        ticStartText,
	}
	xoResult.ParseMode = telebot.ModeMarkdownV2

	xoResult.ReplyMarkup = startInlineKeyboard

	xoResult.SetResultID(string(TicTacToeGameType))
	results = append(results, xoResult)

	// Dot Box result
	dotResult := &telebot.ArticleResult{
		Title:       "نقطه بازی",
		Description: "رو من کلیک کن",
		Text:        dotBoxStartText,
	}
	dotResult.ParseMode = telebot.ModeMarkdownV2

	dotResult.ReplyMarkup = startInlineKeyboard

	dotResult.SetResultID(string(DotBoxGameType))
	results = append(results, dotResult)

	return c.Answer(&telebot.QueryResponse{
		Results:   results,
		CacheTime: 0,
	})
}

func (app *Application) inlineResultHandler(c telebot.Context) error {
	resultId := c.InlineResult().ResultID

	switch GameType(resultId) {
	case TicTacToeGameType:
		return app.ticInlineResultReciever(c, NewTicGame(), ticStartText)
	case DotBoxGameType:
		return app.ticInlineResultReciever(c, NewDotBoxGame(), dotBoxStartText)
	}
	return c.RespondAlert("این بازیرو ندارم!")
}

func (app *Application) ticInlineResultReciever(c telebot.Context, game GameInterface, gameTitle string) error {
	sender := c.InlineResult().Sender
	messageId := c.InlineResult().MessageID

	hub := NewHub(game, c.InlineResult().MessageID)
	hub.Game.SetMessageID(messageId)
	app.Lobby.Hubs.AddHub(hub)
	hub.Game.AddPlayer(NewHumanPlayer(int(sender.ID), sender.FirstName))

	textMessage := gameTitle + "\n\n🕹 بازیکن " + fmt.Sprintf("[%s](tg://user?id=%d)", sender.FirstName, sender.ID) + " منتظر حریفه"

	_, err := c.Bot().Edit(c.InlineResult(), textMessage, app.JoinGameKeyboard(), telebot.ModeMarkdownV2)

	return err
}

func (app *Application) callbackHandler(c telebot.Context) error {
	callback := c.Callback()
	messageId := c.Callback().MessageID

	hub, is_found := app.Lobby.Hubs[messageId]
	if !is_found {
		return c.RespondAlert("این بازی وجود نداره!")
	}

	if callback.Data == "join_hub" {
		return app.JoinHubHandler(c)
	}

	if strings.HasPrefix(callback.Data, string(TicTacToeGameType)+"_") {
		ticGame, ok := hub.Game.(*TicGame)
		if !ok {
			panic("xo callback handler can't convert game interface to ticgame struct")
		}

		return ticGame.TTTCallbackHandlers(c, app.Bot)
	} else if strings.HasPrefix(callback.Data, string(DotBoxGameType)+"_") {
		dotBoxGame, ok := hub.Game.(*DotBoxGame)
		if !ok {
			panic("dot box callback handler can't convert game interface to dot box game struct")
		}

		return dotBoxGame.CallbackHandlers(c, app.Bot)
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
	for _, joinedPlayer := range hub.Game.GetPlayers() {
		if joinedPlayer.TgID == player.TgID {
			return c.RespondText("تو بازی هستی")
		}
	}

	isOk := hub.JoinGame(player, app)
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
