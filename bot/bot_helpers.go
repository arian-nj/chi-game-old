package bot

import (
	"log/slog"

	humanplayer "github.com/arian-nj/chibazi/internals/human_player"
	"gopkg.in/telebot.v4"
)

func ClearMatchmakingMessage(players []*humanplayer.HumanPlayer, bot *telebot.Bot) {
	for _, player := range players {

		go func(inPlayer humanplayer.HumanPlayer) {
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
func SendFoundOpponentMessage(players []*humanplayer.HumanPlayer, bot *telebot.Bot) {
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
