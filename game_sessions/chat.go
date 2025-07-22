package gamesessions

import (
	"fmt"

	humanplayer "github.com/arian-nj/chibazi/internals/human_player"
	"gopkg.in/telebot.v4"
)

type Chat struct {
	IsChatOn bool
}

func (gs *GameSession) HandleChatMessage(bot telebot.API, senderID int, text string) error {
	if !gs.Chat.IsChatOn {
		return nil
	}
	players := gs.GameState.Players()
	var senderPlayer *humanplayer.HumanPlayer
	var recieverPlayer *humanplayer.HumanPlayer

	for _, p := range players {
		if p.TgID == senderID {
			senderPlayer = p
		} else {
			recieverPlayer = p
		}
	}

	_, err := bot.Send(&telebot.User{ID: int64(recieverPlayer.TgID)},
		fmt.Sprintf("*_%s:_* %s", senderPlayer.Name, text), telebot.ModeMarkdownV2)
	return err
}
