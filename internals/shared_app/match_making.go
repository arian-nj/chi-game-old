package sharedapp

import (
	"strconv"
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

func (app *SharedApp) AddTicket(gameType gametype.GameType, newTicket *Ticket) {

	queue := app.MatchMaking.WaitingPlayers[gameType]
	app.MatchMaking.WaitingPlayers[gameType] = append(queue, newTicket)

}

func (app *SharedApp) CheckIsAllowedToPlay(playerId int) bool {
	_, isFound := app.GameSessions.Sessions[strconv.Itoa(playerId)]
	return !isFound
}
