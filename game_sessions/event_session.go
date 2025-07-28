package gamesessions

import "github.com/arian-nj/chibazi/internals/socket"

const (
	ChatType socket.EventType = "chat"
	GameType socket.EventType = "game"
)

type SessionEvent struct {
	Player *SessionPlayer
	Event  *socket.Event
}

func NewSessionEvent(player *SessionPlayer, event *socket.Event) *SessionEvent {
	return &SessionEvent{
		Player: player,
		Event:  event,
	}
}
