package keybul

import (
	"errors"

	gametype "github.com/arian-nj/chibazi/internals/game_type"
	"gopkg.in/telebot.v4"
)

func EditGameMessage(bot telebot.API, msg telebot.Editable, text string, keyboard *telebot.ReplyMarkup) error {
	_, err := bot.Edit(msg, text, keyboard, telebot.ModeMarkdownV2)
	if err != nil && !errors.Is(err, telebot.ErrTrueResult) {
		return err
	}
	return nil
}
func CreateInlineKeyboard(buttonGroups ...[][]telebot.InlineButton) *telebot.ReplyMarkup {
	var inlineKeyboard [][]telebot.InlineButton

	for _, group := range buttonGroups {
		inlineKeyboard = append(inlineKeyboard, group...)
	}

	return &telebot.ReplyMarkup{
		InlineKeyboard: inlineKeyboard,
	}
}

func CreateBotNameInlineButton() [][]telebot.InlineButton {
	return [][]telebot.InlineButton{
		{
			{
				Text: "Chi Bazi | چی بازی",
				URL:  "t.me/" + "ChiBaziBot" + "?start=new",
			},
		},
	}
}

var StartInlineKeyboard = &telebot.ReplyMarkup{
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

func JoinGameKeyboard(game_type gametype.GameType) [][]telebot.InlineButton {
	return [][]telebot.InlineButton{
		{
			{
				Text: "منم بازی",
				Data: string(game_type) + "_join",
			},
		},
	}
}

var EndgameInlineKeyboard = [][]telebot.InlineButton{
	{
		{Text: "🤝 بازی با دوستان", InlineQueryChosenChat: &telebot.SwitchInlineQuery{AllowUserChats: true, AllowGroupChats: true}},
		{Text: "🔄 دوباره", InlineQuery: ""},
	},
}
