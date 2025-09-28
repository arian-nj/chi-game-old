package gamesessions

import (
	"context"
	"log/slog"
	"strconv"

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

type WaitForPlayerCommand struct {
	Session *GameSession
	Creator *SessionPlayer
}

func NewWaitForPlayerCommand(session *GameSession, creatorUser *SessionPlayer) *WaitForPlayerCommand {
	return &WaitForPlayerCommand{
		Session: session,
		Creator: creatorUser,
	}
}

func (wait *WaitForPlayerCommand) Execute() {
}

type JoinSessionCommand struct {
	Session      *GameSession
	JoinedPlayer *SessionPlayer
	AllSession   *AllSession
}

func NewJoinSessionCommand(session *GameSession, JoinedUser *SessionPlayer, allSession *AllSession) *JoinSessionCommand {
	return &JoinSessionCommand{
		Session:      session,
		JoinedPlayer: JoinedUser,
		AllSession:   allSession,
	}
}

func (join *JoinSessionCommand) Execute() {
	join.Session.AddPlayerToSession(join.JoinedPlayer)
	for _, player := range join.Session.Players {
		join.AllSession.Add(strconv.Itoa(player.ID), join.Session)
	}
	join.Session.StartGame()
}

type GameEndedCommand struct {
	Session *GameSession
}

func NewGameEndedSessionCommand(session *GameSession) *GameEndedCommand {
	return &GameEndedCommand{
		Session: session,
	}
}

func (EndGame *GameEndedCommand) Execute() {
}

type GameStartCommand struct {
	Session *GameSession
}

func NewGameStartSessionCommand(session *GameSession) *GameEndedCommand {
	return &GameEndedCommand{
		Session: session,
	}
}

func (EndGame *GameStartCommand) Execute() {
}

type RequestEndSessionCommand struct {
	Session *GameSession
	Player  *SessionPlayer
}

func NewRequestEndSessionCommand(session *GameSession, player *SessionPlayer) *RequestEndSessionCommand {
	return &RequestEndSessionCommand{
		Session: session,
		Player:  player,
	}
}

func (ReqEnd *RequestEndSessionCommand) Execute() {
}
