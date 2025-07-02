package gamebot

import (
	"fmt"

	"gopkg.in/telebot.v4"
)

func CreateInlineKeyboard(buttonGroups ...[][]telebot.InlineButton) *telebot.ReplyMarkup {
	var inlineKeyboard [][]telebot.InlineButton

	for _, group := range buttonGroups {
		inlineKeyboard = append(inlineKeyboard, group...)
	}

	return &telebot.ReplyMarkup{
		InlineKeyboard: inlineKeyboard,
	}
}

func CreateBotNameInlineButton(bot *telebot.Bot) [][]telebot.InlineButton {
	botUsername := bot.Me.Username
	return [][]telebot.InlineButton{
		{
			{
				Text: "Chi Bazi | چی بازی",
				URL:  "t.me/" + botUsername + "?start=new",
			},
		},
	}
}

func CreatePlayersInlineButton(humanPlayers []*HumanPlayer, CurrentPlayerTurn int) [][]telebot.InlineButton {
	buttons := make([][]telebot.InlineButton, 0)
	for index, hplayer := range humanPlayers {

		yourTurn := ""
		if CurrentPlayerTurn == index {
			yourTurn = "🎮"
		}

		emoji := "👤"
		playEmoji := OEmoji
		if index == 0 {
			emoji = "🗿"
			playEmoji = XEmoji
		}

		name := hplayer.Name
		if len(name) > 20 {
			name = name[:20] + "..."
		}

		row := make([]telebot.InlineButton, 2)
		row = append(row, telebot.InlineButton{
			Text: fmt.Sprintf("%s %s (%s) %s", yourTurn, name, playEmoji, emoji),
			URL:  fmt.Sprintf("tg://user?id=%d", hplayer.TgID),
		})
		buttons = append(buttons, row)
	}

	return buttons
}

var startInlineKeyboard = &telebot.ReplyMarkup{
	InlineKeyboard: [][]telebot.InlineButton{
		{
			{
				Text: "واسا بازی رو بسازم",
				Data: "_",
			},
		},
	},
	ResizeKeyboard: true,
}

func (app *Application) JoinGameKeyboard() *telebot.ReplyMarkup {
	var startInlineKeyboard = &telebot.ReplyMarkup{
		InlineKeyboard: [][]telebot.InlineButton{
			{
				{
					Text: "منم بازی",
					Data: "join_hub",
				},
			},
		},
		ResizeKeyboard: true,
	}
	return startInlineKeyboard
}

var EndgameInlineKeyboard = [][]telebot.InlineButton{
	{
		{Text: "🤝 بازی با دوستان", InlineQueryChosenChat: &telebot.SwitchInlineQuery{AllowUserChats: true, AllowGroupChats: true}},
		{Text: "🔄 دوباره", InlineQuery: ""},
	},
}
