package gamesessions

import (
	"fmt"
	"log/slog"

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

	err := SendChatMessageInBot(bot, recieverPlayer.TgID, text, senderPlayer.Name)
	if err != nil {
		slog.Error("can't send chat message from bot", "error", err)
	}
	if recieverPlayer.Socket != nil {
		em := socket.EventMessage(text)
		err := recieverPlayer.Socket.SendEvent(socket.NewEvent(ChatType, em))
		if err != nil {
			slog.Error("can't send chat message to socker", "error", err)
		}
	}

	return nil
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

	chatMessage := newSessionEvent.Event.Data
	if recieverPlayer.Socket != nil {
		return recieverPlayer.Socket.SendEvent(socket.NewEvent(ChatType, chatMessage))
	}

	err := SendChatMessageInBot(gs.Bot, recieverPlayer.TgID, string(chatMessage), senderPlayer.Name)
	if err != nil {
		slog.Error("can't send bot message from got from socket to reciever ", "error", err)
	}

	err = SendChatMessageInBot(gs.Bot, senderPlayer.TgID, string(chatMessage), senderPlayer.Name)
	if err != nil {
		slog.Error("can't send bot message from got from socket to sender", "error", err)
	}
	return nil
}

func SendChatMessageInBot(bot telebot.API, toId int, text string, senderName string) error {
	_, err := bot.Send(&telebot.User{ID: int64(toId)},
		fmt.Sprintf("*_%s:_* %s", senderName, text), telebot.ModeMarkdownV2)
	return err
}
