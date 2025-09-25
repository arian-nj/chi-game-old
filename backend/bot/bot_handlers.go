package bot

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/arian-nj/chibazi/backend/database"
	"github.com/arian-nj/chibazi/backend/games/games"
	keybul "github.com/arian-nj/chibazi/backend/internals/keybul"
	matchmaking "github.com/arian-nj/chibazi/backend/match_making"
	"github.com/jackc/pgx/v5/pgtype"
	"gopkg.in/telebot.v4"
)

func (app *BotApplication) statHandler(c telebot.Context) error {

	if c.Sender().ID != 1909090204 {
		return app.textHandler(c)
	}

	allGamesCount, err := app.Queries.CountSessions(context.Background())
	if err != nil {
		return err
	}

	lastHourCount, err := app.Queries.CountLastHourHub(context.Background())
	if err != nil {
		return err
	}

	lastDayCount, err := app.Queries.CountLastDaySessions(context.Background())
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
	args := c.Args()
	if len(args) > 0 {
		if args[0] == "friends" {
			return app.PlayWithFriendsHandler(c)
		}
	}
	return app.welcomePanel(c)
}

func (app *BotApplication) welcomePanel(c telebot.Context) error {
	return c.Send(
		`خوش اومدید 👋`, keybul.WelcomeReplyKeyboard)
}

func (app *BotApplication) meHandler(c telebot.Context) error {
	personRow, err := app.Queries.GetTgUserByTgID(context.Background(), int(c.Sender().ID))
	if err != nil {
		return err
	}

	text := fmt.Sprintf("id %d \ntg %d", personRow.ID, personRow.TgID)
	return c.Send(text, keybul.WelcomeReplyKeyboard)
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
		return gameSession.BotRequestSendMsg(c.Bot(), senderID, text)
	}

	return c.Send(
		`خوش اومدید 👋`, keybul.WelcomeReplyKeyboard)
}

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
	return app.PlayRandomGameHandler(c, games.XOGameType3X3)
}

func (app *BotApplication) PlayRandomConn4Handler(c telebot.Context) error {
	return app.PlayRandomGameHandler(c, games.Conn4GameType)
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
	return app.PlayRandomGameHandler(c, games.XOGameType5X5)
}

func (app *BotApplication) PlayRandomGameHandler(c telebot.Context, gameType games.GameType) error {
	text := ""
	sender := c.Sender()
	personRow, err := app.Queries.GetTgUserByTgID(context.Background(), int(sender.ID))
	if err != nil {
		slog.Error("can't find user in random game", "error", err)
		return c.Send("پیدات نمیکنم")
	}

	if app.AllSessions.IsSessionEmpty(personRow.ID) == false {
		return c.Send("بازی قبلیت باید تموم بشه")
	}

	if app.MatchMaking.RemovePlayerTicket(personRow.ID) {
		text += "قبلیو لغو کردم\n"
	}

	text += "دنبال حریفم"
	msg, err := c.Bot().Send(sender, text, CancelGameReplyKeyboard)
	if err != nil {
		return err
	}

	newTicket := matchmaking.NewTicket(personRow.Name, personRow.ID, int(sender.ID), gameType)
	app.MatchMaking.PushTicket(newTicket)

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
					Text: keybul.Conn4ButtonText,
				},
				// {
				// 	Text: keybul.Xo5x5ButtonText,
				// },
			},
			{
				{Text: keybul.MainKeyboardButtonText},
			},
		},
		ResizeKeyboard: true,
	}
)

func (app *BotApplication) CancelSearchingForGame(c telebot.Context) error {
	success := app.MatchMaking.RemovePlayerTicket(int(c.Sender().ID))
	if !success {
		slog.Error("can not remove Players Ticket")
	}
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
		gameSession.ShutdownTimer = time.After(0)
		return nil
	}

	gameSession.Chat.IsOn = false
	for _, player := range gameSession.Players {
		_, err := c.Bot().Send(&telebot.User{ID: int64(player.TgID)}, "⛔️چت قطع شد بازی ادامه داره", keybul.WelcomeReplyKeyboard)
		if err != nil {
			slog.Error("can't send close chat message ", "error", err)
		}
	}

	return nil
}
