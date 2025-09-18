package keybul

import (
	"errors"
	"strings"

	"gopkg.in/telebot.v4"
)

func EscapeReserved(text string) string {
	replacer := strings.NewReplacer(
		"_", `\_`,
		"[", `\[`,
		"]", `\]`,
		"(", `\(`,
		")", `\)`,
		"~", `\~`,
		"`", "\\`", // The backslash needs to be escaped in a Go string literal
		">", `\>`,
		"#", `\#`,
		"+", `\+`,
		"-", `\-`,
		"=", `\=`,
		"|", `\|`,
		"{", `\{`,
		"}", `\}`,
		".", `\.`,
		"!", `\!`,
	)
	return replacer.Replace(text)
}

func EditMessage(bot telebot.API, msg telebot.Editable, text string, keyboard *telebot.ReplyMarkup) error {
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

func ContinueInWebButton() [][]telebot.InlineButton {
	return [][]telebot.InlineButton{
		{
			{
				Text: "🎮 ادامه بازی در وب",
				URL:  "t.me/ChiGameBot/game",
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

// func JoinGameKeyboard() [][]telebot.InlineButton {
var JoinGameInlineButtons = [][]telebot.InlineButton{
	{
		{
			Text: "منم بازی",
			Data: "join",
		},
	},
}

// }
func EndGameInlineKeyboard(isVia bool) [][]telebot.InlineButton {
	if !isVia {
		return [][]telebot.InlineButton{}
	}
	var EndgameInlineKeyboard = [][]telebot.InlineButton{
		{
			{Text: "🤝 بازی با دوستان", InlineQueryChosenChat: &telebot.SwitchInlineQuery{AllowUserChats: true, AllowGroupChats: true}},
			{Text: "🔄 دوباره", InlineQuery: ""},
		},
	}
	return EndgameInlineKeyboard
}

var (
	WelcomeReplyKeyboard = &telebot.ReplyMarkup{
		ReplyKeyboard: [][]telebot.ReplyButton{
			{
				{Text: PlayWithFriendsButtonText},
				{Text: PlayWithRandomPlayerText},
			},
		},
		ResizeKeyboard: true,
	}
)
