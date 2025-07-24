package keybul

import (
	"gopkg.in/telebot.v4"
)

var StopChatReplyKeyboard = &telebot.ReplyMarkup{
	ReplyKeyboard: [][]telebot.ReplyButton{
		{
			{Text: StopChatButtonText},
		},
	},
	ResizeKeyboard: true,
}
