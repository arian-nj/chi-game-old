package bot

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/arian-nj/chibazi/backend/database"
	gamesessions "github.com/arian-nj/chibazi/backend/game_sessions"
	"github.com/arian-nj/chibazi/backend/games/conn4"
	"github.com/arian-nj/chibazi/backend/games/games"
	"github.com/arian-nj/chibazi/backend/games/xo"
	"github.com/arian-nj/chibazi/backend/internals/keybul"
	"gopkg.in/telebot.v4"
)

func (app *BotApplication) inlineResultFeedbackHandler(c telebot.Context) error {
	resultId := c.InlineResult().ResultID
	viaMessageID := c.InlineResult().MessageID

	newSessionRow, err := app.Queries.CreateSession(context.Background(), string(gamesessions.PrivateSession))
	if err != nil {
		slog.Error("can't create session inline feedback")
		return err
	}

	newSession := gamesessions.NewGameSession(app.Bot, app.Queries, newSessionRow.ID, app.AllSessions)
	sessionTgListen := gamesessions.NewSessionTelegramViaListener(app.Bot, viaMessageID)
	newSession.Subscribe(sessionTgListen)

	var newGame games.Game
	gameType := games.GameType(resultId)
	switch gameType {
	case games.XOGameType3X3:
		newGame = xo.NewXOGame(newSession.SessionCtx, games.XOGameType3X3, app.Queries)
	case games.XOGameType5X5:
		newGame = xo.NewXOGame(newSession.SessionCtx, games.XOGameType5X5, app.Queries)
	case games.Conn4GameType:
		newGame = conn4.NewConn4State(newSession.SessionCtx, games.XOGameType5X5, app.Queries)
	default:
		return c.RespondAlert("این بازیرو ندارم!")
	}

	if newGame == nil {
		err := fmt.Errorf("new game is nil")
		slog.Error("inline feedback handler", "error", err)
		return err
	}

	newSession.GameState = newGame
	app.AllSessions.Add(viaMessageID, newSession)

	personRow, err := app.Queries.GetTgUserByTgID(context.Background(), int(c.Sender().ID))
	if err != nil {
		slog.Error("can not get creator user in feedback")
		return err
	}

	creatorPlayer := gamesessions.NewSessionPlayer(personRow.ID, personRow.TgID, personRow.Name)
	newSession.AddPlayerToSession(creatorPlayer)

	newSession.RunBgMonitor()

	newSession.PushCommand(
		gamesessions.NewWaitForPlayerCommand(newSession, creatorPlayer),
	)

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
		Text:        xo.XOStartText,
	}
	xoResult3x3.ParseMode = telebot.ModeMarkdownV2
	xoResult3x3.ReplyMarkup = keybul.StartInlineKeyboard
	xoResult3x3.SetResultID(string(games.XOGameType3X3))
	results = append(results, xoResult3x3)

	conn4Result := &telebot.ArticleResult{
		Title:       "دوز  سنگی",
		Description: "رو من کلیک کن",
		Text:        conn4.Conn4StartText,
	}
	conn4Result.ParseMode = telebot.ModeMarkdownV2
	conn4Result.ReplyMarkup = keybul.StartInlineKeyboard
	conn4Result.SetResultID(string(games.Conn4GameType))
	results = append(results, conn4Result)

	return c.Answer(&telebot.QueryResponse{
		Results:   results,
		CacheTime: 0,
	})
}
