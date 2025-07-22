package bot

import (
	"log/slog"

	consoleplayer "github.com/arian-nj/chibazi/internals/console_player"
	"gopkg.in/telebot.v4"
)

func ClearMatchmakingMessage(players []*consoleplayer.ConsolePlayer, bot *telebot.Bot) {
	for _, player := range players {

		go func(inPlayer consoleplayer.ConsolePlayer) {
			err := bot.Delete(&inPlayer)
			if err != nil {
				slog.Error("can't delete match making message", "err", err)
				return
			}
		}(*player)

		msg, err := bot.Send(&telebot.User{ID: int64(player.TgID)}, "بازی")
		if err != nil {
			slog.Error("can't send base game message", "err", err)
		}
		player.MessageID = msg.ID
	}
}

// SendFoundOpponentMessage sends a message to all players that they have found an opponent
func SendFoundOpponentMessage(players []*consoleplayer.ConsolePlayer, bot *telebot.Bot) {
	for _, player := range players {
		for _, oppPlayer := range players {
			if player.TgID == oppPlayer.TgID {
				continue
			}
			_, err := bot.Send(&telebot.User{ID: int64(player.TgID)}, FoundOpponentText(oppPlayer.Name), telebot.ModeMarkdownV2, StopChatReplyKeyboard)
			if err != nil {
				slog.Error("can't send found opponent message ", "error", err)
			}
		}
	}
}
