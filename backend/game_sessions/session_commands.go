package gamesessions

import (
	"context"
	"log/slog"

	"github.com/arian-nj/chibazi/backend/database"
	"github.com/arian-nj/chibazi/backend/internals/utils"
)

type MessageCommand struct {
	Text     string
	Sender   *SessionPlayer
	Reciever *SessionPlayer
	Session  *GameSession
}

func NewMessageCommand(session *GameSession, text string, senderPlayer *SessionPlayer, recPlayer *SessionPlayer) *MessageCommand {
	return &MessageCommand{
		Text:     text,
		Session:  session,
		Sender:   senderPlayer,
		Reciever: recPlayer,
	}
}

func (message *MessageCommand) Execute() {
	session := message.Session

	utils.RunBackgroundTask(func() {
		_, err := session.Queries.CreateSessionMessage(context.Background(), database.CreateSessionMessageParams{
			SessionID: session.ID,
			PlayerID:  message.Sender.ID,
			Message:   message.Text,
		})
		if err != nil {
			slog.Error("error creating new message in db")
		}
	})
}
