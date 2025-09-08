package gamesessions

import (
	"log/slog"

	sessionv1 "github.com/arian-nj/chibazi/backend/gen/session/v1"
	"github.com/arian-nj/chibazi/backend/internals/commander"
)

type SessionSocketListener struct {
	UserID int
}

func NewSessionSocketListener(userID int) *SessionSocketListener {
	return &SessionSocketListener{
		UserID: userID,
	}
}
func (sl *SessionSocketListener) Update(command commander.Command) {
	switch c := command.(type) {
	case *MessageCommand:
		if c.Reciever.ID == sl.UserID {
			err := SendChatMessageInWeb(c.Reciever, c.Sender, c.Text)
			if err != nil {
				slog.Error("Can not send chat message in session socket", "error", err)
			}
		}
	}
}

func (session *GameSession) SocketRequestSendMsg(sessionPlayer *SessionPlayer, chatMsgReq *sessionv1.ChatMessageRequest) {
	if !session.Chat.IsOn {
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

	for _, p := range session.Players {
		if p.TgID == senderID {
			senderPlayer = p
		} else {
			recieverPlayer = p
		}
	}

	session.PushCommand(NewMessageCommand(session, messageText, senderPlayer, recieverPlayer))
}
