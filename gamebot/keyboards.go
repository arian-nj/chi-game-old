package gamebot

import (
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
