package gamebot

import (
	"context"
	"fmt"

	"github.com/arian-nj/chibazi/games/dotbox_console"
	xoconsole "github.com/arian-nj/chibazi/games/xo_console"
	gametype "github.com/arian-nj/chibazi/internals/game_type"
	keybul "github.com/arian-nj/chibazi/internals/keybul"
	"gopkg.in/telebot.v4"
)

func (app *Application) inlineQueryHandler(c telebot.Context) error {
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

	// Dot Box result
	dotResult := &telebot.ArticleResult{
		Title:       "نقطه بازی",
		Description: "رو من کلیک کن",
		Text:        dotbox_console.DotBoxStartText,
	}
	dotResult.ParseMode = telebot.ModeMarkdownV2

	dotResult.ReplyMarkup = keybul.StartInlineKeyboard

	dotResult.SetResultID(string(gametype.DotBoxGameType))
	results = append(results, dotResult)

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

func (app *Application) statHandler(c telebot.Context) error {
	count, err := app.Queries.CountHubs(context.Background())
	if err != nil {
		return err
	}

	c.Send(fmt.Sprintf("تعداد بازی ها: %d", count))
	return nil
}

var (
	welcomeInlinequery = &telebot.ReplyMarkup{
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

func (app *Application) welcomeHandler(c telebot.Context) error {
	return c.Send(
		`خوش اومدید 👋
دکمه بازی با دوستان رو بزن تا تو هر چت یا گروهی با دوستات بازی کنی
	`, welcomeInlinequery)
}
