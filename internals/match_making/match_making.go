package matchmaking

import (
	"sync"
	"time"

	gametype "github.com/arian-nj/chibazi/internals/game_type"
)

type MatchMaking struct {
	WaitingPlayers map[gametype.GameType][]*Ticket
	Mutex          sync.Mutex
}

type Ticket struct {
	Name      string
	UserID    int
	MessageID int
	GameType  gametype.GameType
	Timestamp time.Time
}

func NewTicket(name string, userID, messageID int, gameType gametype.GameType) *Ticket {
	return &Ticket{
		UserID:    userID,
		Name:      name,
		MessageID: messageID,
		GameType:  gameType,
	}
}
