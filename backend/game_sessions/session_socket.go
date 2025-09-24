package gamesessions

import (
	"fmt"
	"log/slog"

	"github.com/arian-nj/chibazi/backend/games/games"
	sessionv1 "github.com/arian-nj/chibazi/backend/gen/session/v1"
	"github.com/arian-nj/chibazi/backend/internals/commander"
)

type SessionSocketListener struct {
	Player *SessionPlayer
}

func NewSessionSocketListener(player *SessionPlayer) *SessionSocketListener {
	return &SessionSocketListener{
		Player: player,
	}
}
func (sl *SessionSocketListener) Update(command commander.Command) {
	switch c := command.(type) {
	case *MessageCommand:
		if c.Reciever.ID == sl.Player.ID {
			err := SendChatMessageInWeb(c.Reciever, c.Sender, c.Text)
			if err != nil {
				slog.Error("Can not send chat message in session socket", "error", err)
			}
		}
	case *GameStartCommand:
		SendGametypeOverSocket(c.Session, sl.Player)
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

func SendGametypeOverSocket(session *GameSession, player *SessionPlayer) {
	if session.GameState == nil {
		slog.Error("can not send game type message game state is nil")
		return
	}
	gameData := session.GameState.GetGameData()
	currentGameType := gameData.GameType

	var protoGameType sessionv1.GameType
	switch currentGameType {
	case games.XOGameType3X3:
		protoGameType = sessionv1.GameType_GAME_TYPE_XO3X3
	case games.Conn4GameType:
		protoGameType = sessionv1.GameType_GAME_TYPE_CONN4
	default:
		slog.Error("unknown game type to send", "game_type", currentGameType)
		return
	}

	message := sessionv1.SessionMessage{Content: &sessionv1.SessionMessage_GameType{
		GameType: &sessionv1.ChangeGameTypeMessage{
			GameType: protoGameType,
		},
	}}

	player.Socket.SendMessage(&message)
}
