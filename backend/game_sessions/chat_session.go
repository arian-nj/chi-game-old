package gamesessions

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/arian-nj/chibazi/backend/database"
	sessionv1 "github.com/arian-nj/chibazi/backend/gen/session/v1"
	"github.com/arian-nj/chibazi/backend/internals/socket"
	"github.com/arian-nj/chibazi/backend/internals/utils"
	"gopkg.in/telebot.v4"
)

type Chat struct {
	IsOn bool
}

func (gs *GameSession) HandleBotChatMessage(bot telebot.API, senderID int, messageText string) error {
	if !gs.Chat.IsOn {
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

	err := SendChatMessageInBot(bot, recieverPlayer.TgID, messageText, senderPlayer.Name)
	if err != nil {
		slog.Error("can't send chat message from bot", "error", err)
	}

	if recieverPlayer.Socket != nil {
		err := recieverPlayer.Socket.SendNewEvent(socket.ChatEventType, messageText)
		if err != nil {
			slog.Error("can't send chat message to socket", "error", err, "message", messageText)
		}
	}

	utils.RunBackgroundTask(func() {
		_, err := gs.Queries.CreateSessionMessage(context.Background(), database.CreateSessionMessageParams{
			SessionID: gs.ID,
			PlayerID:  senderPlayer.ID,
			Message:   messageText,
		})
		if err != nil {
			slog.Error("error creating new message in db")
		}
	})

	return nil
}

func (gs *GameSession) HandleWebChatMessage(sessionPlayer *SessionPlayer, chatMsgReq *sessionv1.ChatMessageRequest) {
	if !gs.Chat.IsOn {
		return
	}

	messageText := chatMsgReq.Text
	if len(messageText) > 256 {
		slog.Error("message is to long")
		return
	}
	senderID := sessionPlayer.TgID

	var senderPlayer *SessionPlayer
	var recieverPlayer *SessionPlayer

	for _, p := range gs.Players {
		if p.TgID == senderID {
			senderPlayer = p
		} else {
			recieverPlayer = p
		}
	}

	if recieverPlayer.Socket != nil {
		newChatMsg := &sessionv1.SessionMessage{
			Content: &sessionv1.SessionMessage_Chat{
				Chat: &sessionv1.ChatMessage{
					PlayerId: int32(senderPlayer.ID),
					Text:     messageText,
				},
			},
		}
		err := recieverPlayer.Socket.SendMessage(newChatMsg)
		if err != nil {
			slog.Error("error seding message to socket", "error", err)
		}
	}

	err := SendChatMessageInBot(gs.Bot, recieverPlayer.TgID, string(messageText), senderPlayer.Name)
	if err != nil {
		slog.Error("can't send bot message from got from socket to reciever ", "error", err)
	}

	err = SendChatMessageInBot(gs.Bot, senderPlayer.TgID, string(messageText), senderPlayer.Name)
	if err != nil {
		slog.Error("can't send bot message from got from socket to sender", "error", err)
	}

	utils.RunBackgroundTask(func() {
		_, err := gs.Queries.CreateSessionMessage(context.Background(), database.CreateSessionMessageParams{
			SessionID: gs.ID,
			PlayerID:  senderPlayer.ID,
			Message:   string(messageText),
		})
		if err != nil {
			slog.Error("error creating new message in db")
		}
	})
}

func SendChatMessageInBot(bot telebot.API, toId int, text string, senderName string) error {
	_, err := bot.Send(&telebot.User{ID: int64(toId)},
		fmt.Sprintf("*_%s:_* %s", senderName, text), telebot.ModeMarkdownV2)
	return err
}
