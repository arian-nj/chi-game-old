package bot

import (
	"context"

	gamesessions "github.com/arian-nj/chibazi/game_sessions"
	xoconsole "github.com/arian-nj/chibazi/games/xo_console"
	gametype "github.com/arian-nj/chibazi/internals/game_type"
	"github.com/arian-nj/chibazi/internals/keybul"
	"gopkg.in/telebot.v4"
)

func (app *BotApplication) inlineResultFeedbackHandler(c telebot.Context) error {
	resultId := c.InlineResult().ResultID
	messageID := c.InlineResult().MessageID

	newSession, err := app.Queries.CreateSession(context.Background())
	if err != nil {
		return err
	}

	switch gametype.GameType(resultId) {
	case gametype.XOGameType3X3:
		newXOGame := xoconsole.NewXOGame(gametype.XOGameType3X3, app.Queries)
		newXOGame.ViaMessageId = messageID

		newGameSession := gamesessions.NewGameSession(app.AllSessions, app.Bot, gametype.XOGameType3X3, newXOGame, newSession.ID)

		app.AllSessions.Add(messageID, newGameSession)

		return newXOGame.SendJoinPanelAddSender(c)

	case gametype.XOGameType5X5:
		newXOGame := xoconsole.NewXOGame(gametype.XOGameType5X5, app.Queries)
		newXOGame.ViaMessageId = messageID

		newGameSession := gamesessions.NewGameSession(app.AllSessions, app.Bot, gametype.XOGameType3X3, newXOGame, newSession.ID)

		app.AllSessions.Add(messageID, newGameSession)

		return newXOGame.SendJoinPanelAddSender(c)
	}
	return c.RespondAlert("این بازیرو ندارم!")
}

func (app *BotApplication) inlineQueryHandler(c telebot.Context) error {
	results := telebot.Results{}
	// XO result
	xoResult3x3 := &telebot.ArticleResult{
		Title:       "دوز بازی ۳ در ۳",
		Description: "رو من کلیک کن",
		Text:        xoconsole.XOStartText,
	}
	xoResult3x3.ParseMode = telebot.ModeMarkdownV2

	xoResult3x3.ReplyMarkup = keybul.StartInlineKeyboard

	xoResult3x3.SetResultID(string(gametype.XOGameType3X3))
	results = append(results, xoResult3x3)

	xoResult5x5 := &telebot.ArticleResult{
		Title:       "دوز بازی  ۵ در ۵",
		Description: "رو من کلیک کن",
		Text:        xoconsole.XOStartText,
	}
	xoResult5x5.ParseMode = telebot.ModeMarkdownV2

	xoResult5x5.ReplyMarkup = keybul.StartInlineKeyboard

	xoResult5x5.SetResultID(string(gametype.XOGameType5X5))
	results = append(results, xoResult5x5)

	// // Dot Box result
	// dotResult := &telebot.ArticleResult{
	// 	Title:       "نقطه بازی",
	// 	Description: "رو من کلیک کن",
	// 	Text:        dotbox_console.DotBoxStartText,
	// }
	// dotResult.ParseMode = telebot.ModeMarkdownV2
	//
	// dotResult.ReplyMarkup = keybul.StartInlineKeyboard
	//
	// dotResult.SetResultID(string(gametype.DotBoxGameType))
	// results = append(results, dotResult)
	//
	// // Web Dot Box result
	// webDotResult := &telebot.ArticleResult{
	// 	Title:       "نقطه بازی گرافیکی",
	// 	Description: "رو من کلیک کن",
	// 	Text:        consolegames.DotBoxStartText,
	// }
	// webDotResult.ParseMode = telebot.ModeMarkdownV2
	//
	// webDotResult.ReplyMarkup = keybul.StartInlineKeyboard
	//
	// webDotResult.SetResultID(string(gametype.WebDotBoxGameType))
	// results = append(results, webDotResult)
	//
	return c.Answer(&telebot.QueryResponse{
		Results:   results,
		CacheTime: 0,
	})
}
