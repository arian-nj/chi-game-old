package gamesessions

import (
	sessionv1 "github.com/arian-nj/chibazi/backend/gen/session/v1"
)

type SessionEvent struct {
	Player *SessionPlayer
	Event  *sessionv1.SessionMessage
}

func NewSessionEvent(player *SessionPlayer, event *sessionv1.SessionMessage) *SessionEvent {
	return &SessionEvent{
		Player: player,
		Event:  event,
	}
}
