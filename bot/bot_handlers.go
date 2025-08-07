package bot

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/arian-nj/chibazi/database"
	xoconsole "github.com/arian-nj/chibazi/games/xo_console"
	gametype "github.com/arian-nj/chibazi/internals/game_type"
	keybul "github.com/arian-nj/chibazi/internals/keybul"
	matchmaking "github.com/arian-nj/chibazi/match_making"
	"github.com/jackc/pgx/v5/pgtype"
	"gopkg.in/telebot.v4"
)

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

func (app *BotApplication) statHandler(c telebot.Context) error {
	if c.Sender().ID != 1909090204 {
		return app.textHandler(c)
	}

	allGamesCount, err := app.Queries.CountGameSessions(context.Background())
	if err != nil {
		return err
	}

	lastHourCount, err := app.Queries.CountLastHourHub(context.Background())
	if err != nil {
		return err
	}

	lastDayCount, err := app.Queries.CountLastDayHubs(context.Background())
	if err != nil {
		return err
	}

	allUsers, err := app.Queries.CountAllTgUsers(context.Background())
	if err != nil {
		return err
	}
	now := time.Now()
	newUsersLastHour, err := app.Queries.CountUsersTgCreatedBetween(context.Background(), database.CountUsersTgCreatedBetweenParams{
		CreatedAt:   pgtype.Timestamp{Time: now.Add(-time.Hour)},
		CreatedAt_2: pgtype.Timestamp{Time: now},
	})

	newUsersLastDay, err := app.Queries.CountUsersTgCreatedBetween(context.Background(), database.CountUsersTgCreatedBetweenParams{
		CreatedAt:   pgtype.Timestamp{Time: now.Add(-(time.Duration(now.Hour()) + time.Duration(now.Minute()) + time.Duration(now.Second())))},
		CreatedAt_2: pgtype.Timestamp{Time: now},
	})

	text := ""
	text += "*All Games Count*\n" +
		fmt.Sprintf("ever played games: %d\n", allGamesCount) +
		fmt.Sprintf("last day games: %d\n", lastDayCount) +
		fmt.Sprintf("last hour games: %d\n", lastHourCount) +
		"\n\n" +

		fmt.Sprintf("all users: %d\n", allUsers) +
		fmt.Sprintf("new user last hour: %d\n", newUsersLastHour) +
		fmt.Sprintf("last day users: %d\n", newUsersLastDay)

	return c.Send(text, telebot.ModeMarkdownV2)
}

func (app *BotApplication) welcomeHandler(c telebot.Context) error {
	return c.Send(
		`خوش اومدید 👋`, welcomeReplyKeyboard)
}
func (app *BotApplication) textHandler(c telebot.Context) error {
	if !c.Message().Private() {
		return nil
	}

	text := c.Text()
	senderID := int(c.Sender().ID)

	// ensure game is in not Via Message
	gameSession, ok := app.AllSessions.Get(strconv.Itoa(senderID))
	if ok {
		return gameSession.HandleBotChatMessage(c.Bot(), senderID, text)
	}

	return c.Send(
		`خوش اومدید 👋`, welcomeReplyKeyboard)
}

var (
	welcomeReplyKeyboard = &telebot.ReplyMarkup{
		ReplyKeyboard: [][]telebot.ReplyButton{
			{
				{Text: keybul.PlayWithFriendsButtonText},
				{Text: keybul.PlayWithRandomPlayerText},
			},
		},
		ResizeKeyboard: true,
	}
)

func (app *BotApplication) PlayWithFriendsHandler(c telebot.Context) error {
	return c.Send(
		`دکمه بازی با دوستان رو بزن تا تو هر چت یا گروهی با دوستات بازی کنی`, playWithFriendsInline)
}

var (
	playWithFriendsInline = &telebot.ReplyMarkup{
		InlineKeyboard: [][]telebot.InlineButton{
			{
				{
					Text:                  "بازی با دوستان",
					InlineQueryChosenChat: &telebot.SwitchInlineQuery{AllowUserChats: true, AllowGroupChats: true},
				},
			},
		},
	}
)

func (app *BotApplication) PlayWithRandomPlayerHandler(c telebot.Context) error {
	return c.Send("چی بازی؟", WhatRandomGameReplyKeyboard)
}

func (app *BotApplication) PlayRandomXO3X3Handler(c telebot.Context) error {
	return app.PlayRandomGameHandler(c, gametype.XOGameType3X3)
}

var (
	CancelGameReplyKeyboard = &telebot.ReplyMarkup{
		ReplyKeyboard: [][]telebot.ReplyButton{
			{
				{Text: keybul.CancelGameButtonText},
			},
		},
		ResizeKeyboard: true,
	}
)

func (app *BotApplication) PlayRandomXO5X5Handler(c telebot.Context) error {
	return app.PlayRandomGameHandler(c, gametype.XOGameType5X5)
}

func (app *BotApplication) PlayRandomGameHandler(c telebot.Context, gameType gametype.GameType) error {
	text := ""
	sender := c.Sender()

	if app.AllSessions.IsSessionEmpty(int(sender.ID)) == false {
		return c.Send("بازی قبلیت باید تموم بشه")
	}

	if app.MatchMaking.RemovePlayerFromMatchMaking(int(sender.ID)) {
		text += "قبلیو لغو کردم\n"
	}

	text += "دنبال حریفم"
	msg, err := c.Bot().Send(sender, text, CancelGameReplyKeyboard)
	if err != nil {
		return err
	}
	newTicket := matchmaking.NewTicket(sender.FirstName, int(sender.ID), gameType)
	app.MatchMaking.AddTicket(newTicket)

	select {
	case <-newTicket.MatchFoundChan:
		return c.Bot().Delete(msg)
	case <-time.After(90 * time.Second):
		_, err = c.Bot().Edit(msg, "حریف پیدا نکردم")
	}
	return err
}

var (
	WhatRandomGameReplyKeyboard = &telebot.ReplyMarkup{
		ReplyKeyboard: [][]telebot.ReplyButton{
			{
				{
					Text: keybul.Xo3x3ButtonText,
				},
				{
					Text: keybul.Xo5x5ButtonText,
				},
			},
			{
				{Text: keybul.MainKeyboardButtonText},
			},
		},
		ResizeKeyboard: true,
	}
)

func (app *BotApplication) CancelSearchingForGame(c telebot.Context) error {
	app.MatchMaking.Mutex.Lock()
	defer app.MatchMaking.Mutex.Unlock()

	app.MatchMaking.RemovePlayerFromMatchMaking(int(c.Sender().ID))
	return c.Send("لغوش کردم 😔", WhatRandomGameReplyKeyboard)
}

func (app *BotApplication) StopChatHandler(c telebot.Context) error {
	sender := c.Sender()

	app.AllSessions.Mutex.Lock()
	gameSession, isFound := app.AllSessions.Sessions[strconv.Itoa(int(sender.ID))]
	app.AllSessions.Mutex.Unlock()

	if !isFound {
		return c.RespondText("بازی فعالی نداری")
	}

	if gameSession.IsGameEnded {
		gameSession.CleanAndDisconnect(app.AllSessions)
		return nil
	}

	gameSession.Chat.IsChatOn = false
	for _, player := range gameSession.GameState.Players() {
		_, err := c.Bot().Send(&telebot.User{ID: int64(player.TgID)}, "⛔️چت قطع شد بازی ادامه داره", welcomeReplyKeyboard)
		if err != nil {
			slog.Error("can't send close chat message ", "error", err)
		}
	}

	return nil
}
