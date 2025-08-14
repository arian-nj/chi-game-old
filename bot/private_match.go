package bot

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/arian-nj/chibazi/database"
	gamesessions "github.com/arian-nj/chibazi/game_sessions"
	xoconsole "github.com/arian-nj/chibazi/games/xo_console"
	gametype "github.com/arian-nj/chibazi/internals/game_type"
	"github.com/arian-nj/chibazi/internals/keybul"
	"gopkg.in/telebot.v4"
)

func (app *BotApplication) inlineResultFeedbackHandler(c telebot.Context) error {
	resultId := c.InlineResult().ResultID
	messageID := c.InlineResult().MessageID

	newSessionRow, err := app.Queries.CreateSession(context.Background(), string(gamesessions.RandomSession))
	if err != nil {
		return err
	}

	newGameSession := gamesessions.NewGameSession(app.Bot, app.Queries, newSessionRow.ID)

	var newGame gamesessions.Game
	gameType := gametype.GameType(resultId)
	switch gameType {
	case gametype.XOGameType3X3:
		newXOGame := xoconsole.NewXOGame(newGameSession.SessionCtx, gametype.XOGameType3X3, app.Bot, app.Queries)
		newXOGame.ViaMessageId = messageID
		newGame = newXOGame
	case gametype.XOGameType5X5:
		newXOGame := xoconsole.NewXOGame(newGameSession.SessionCtx, gametype.XOGameType5X5, app.Bot, app.Queries)
		newXOGame.ViaMessageId = messageID
		newGame = newXOGame
	default:
		return c.RespondAlert("این بازیرو ندارم!")
	}

	if newGame == nil {
		err := fmt.Errorf("new game is nil")
		slog.Error("inline feedback handler", "error", err)
		return err
	}

	newGameSession.GameState = newGame
	newGameSession.RunBgTask(app.AllSessions)

	app.AllSessions.Add(messageID, newGameSession)

	err = newGame.SendJoinPanelAddSender(c)
	if err != nil {
		slog.Error("can't send join panel", "error", err)
		return err
	}
	_, err = app.Queries.CreateSessionGame(context.Background(), database.CreateSessionGameParams{
		SessionID: newSessionRow.ID,
		GameType:  string(gameType),
	})

	if err != nil {
		slog.Error("error creating session game in when creating private match", "error", err)
		return err
	}
	return nil
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
