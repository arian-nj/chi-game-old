package gamesessions

import (
	"log/slog"

	"github.com/arian-nj/chibazi/backend/internals/commander"
	"gopkg.in/telebot.v4"
)

type SessionTelegramListener struct {
	Bot          *telebot.Bot
	ViaMessageId string // Via Bots
	UserID       int
	TgID         int
}

func NewSessionTelegramListener(playerID int, TgID int, bot *telebot.Bot, viaMessageID string) *SessionTelegramListener {
	return &SessionTelegramListener{
		Bot:          bot,
		ViaMessageId: viaMessageID,
		UserID:       playerID,
		TgID:         TgID,
	}
}

func (tg *SessionTelegramListener) Update(command commander.Command) {
	switch c := command.(type) {
	case *MessageCommand:
		if c.Reciever.ID == tg.UserID {
			err := SendChatMessageInBot(tg.Bot, tg.TgID, c.Text, c.Sender.Name)
			if err != nil {
				slog.Error("Can not send chat message in session telegram", "error", err)
			}
		}
	}
}

// Handler
func (session *GameSession) BotRequestSendMsg(bot telebot.API, senderID int, messageText string) error {
	if !session.Chat.IsOn {
		return nil
	}
	if len(messageText) > 256 {
		slog.Error("message is to long")
		return nil
	}

	var senderPlayer *SessionPlayer
	var recieverPlayer *SessionPlayer

	for _, p := range session.Players {
		if p.TgID == senderID {
			senderPlayer = p
		} else {
			recieverPlayer = p
		}
	}
	session.PushCommand(NewMessageCommand(session, messageText, senderPlayer, recieverPlayer))
	return nil
}
