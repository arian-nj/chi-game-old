package gamesessions

import "github.com/arian-nj/chibazi/internals/socket"

type SessionEvent struct {
	Player *SessionPlayer
	Event  *socket.SocketEvent
}

func NewSessionEvent(player *SessionPlayer, event *socket.SocketEvent) *SessionEvent {
	return &SessionEvent{
		Player: player,
		Event:  event,
	}
}
