package gamesessions

import (
	"fmt"

	"github.com/arian-nj/chibazi/internals/socket"
	"gopkg.in/telebot.v4"
)

type Chat struct {
	IsChatOn bool
}

func (gs *GameSession) HandleBotChatMessage(bot telebot.API, senderID int, text string) error {
	if !gs.Chat.IsChatOn {
		return nil
	}
	var senderPlayer *SessionPlayer
	var recieverPlayer *SessionPlayer

	for _, p := range gs.Players {
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

func (gs *GameSession) HandleWebChatMessage(newSessionEvent *SessionEvent) error {
	if !gs.Chat.IsChatOn {
		return nil
	}
	senderID := newSessionEvent.Player.TgID

	var senderPlayer *SessionPlayer
	var recieverPlayer *SessionPlayer

	for _, p := range gs.Players {
		if p.TgID == senderID {
			senderPlayer = p
		} else {
			recieverPlayer = p
		}
	}
	_ = senderPlayer
	return recieverPlayer.Socket.SendEvent(socket.NewEvent(ChatType, newSessionEvent.Event.Data))
}
