package bot

import (
	"log/slog"

	humanplayer "github.com/arian-nj/chibazi/internals/human_player"
	"github.com/arian-nj/chibazi/internals/keybul"
	"gopkg.in/telebot.v4"
)

// SendFoundOpponentMessage sends a message to all players that they have found an opponent
func SendFoundOpponentMessage(players []*humanplayer.HumanPlayer, bot *telebot.Bot) {
	for _, player := range players {
		for _, oppPlayer := range players {
			if player.TgID == oppPlayer.TgID {
				continue
			}
			_, err := bot.Send(&telebot.User{ID: int64(player.TgID)}, FoundOpponentText(oppPlayer.Name), telebot.ModeMarkdownV2, keybul.StopChatReplyKeyboard)
			if err != nil {
				slog.Error("can't send found opponent message ", "error", err)
			}
		}
	}
}
