package gamesessions

import (
	"fmt"
	"log/slog"

	"github.com/arian-nj/chibazi/internals/keybul"
	"gopkg.in/telebot.v4"
)

// SendFoundOpponentMessage sends a message to all players that they have found an opponent
func SendFoundOpponentMessage(players []*SessionPlayer, bot *telebot.Bot) {
	for _, player := range players {
		for _, oppPlayer := range players {
			if player.TgID == oppPlayer.TgID {
				continue
			}
			_, err := bot.Send(&telebot.User{ID: int64(player.TgID)}, FoundOpponentText(oppPlayer.Name), telebot.ModeMarkdownV2, keybul.StopChatReplyKeyboard)
			if err != nil {
				slog.Error("can't send found opponent message ", "error", err, "text", FoundOpponentText(oppPlayer.Name))

			}
		}
	}
}

func FoundOpponentText(oppName string) string {
	text := ""
	text += "🕹 بازی شروع شد☝️\n\n"
	text += fmt.Sprintf("👀 به حریفت *%s* سلام کن 🤝", oppName)
	return text
}
