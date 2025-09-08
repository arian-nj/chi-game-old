package gamesessions

import (
	"fmt"

	sessionv1 "github.com/arian-nj/chibazi/backend/gen/session/v1"
	"gopkg.in/telebot.v4"
)

type Chat struct {
	IsOn bool
}

func SendChatMessageInBot(bot telebot.API, toId int, text string, senderName string) error {
	_, err := bot.Send(&telebot.User{ID: int64(toId)},
		fmt.Sprintf("*_%s:_* %s", senderName, text), telebot.ModeMarkdownV2)
	return err
}

func SendChatMessageInWeb(recieverPlayer *SessionPlayer, senderPlayer *SessionPlayer, messageText string) error {
	if recieverPlayer.Socket != nil {
		newChatMsg := &sessionv1.SessionMessage{
			Content: &sessionv1.SessionMessage_Chat{
				Chat: &sessionv1.ChatMessage{
					PlayerId: int64(senderPlayer.ID),
					Text:     messageText,
				},
			},
		}
		return recieverPlayer.Socket.SendMessage(newChatMsg)
	}
	return fmt.Errorf("no socket found")
}
