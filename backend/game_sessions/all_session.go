package gamesessions

import (
	"strconv"
	"sync"
)

type AllSession struct {
	Sessions map[string]*GameSession
	Mutex    sync.Mutex
}

func (all *AllSession) Get(look_for string) (*GameSession, bool) {
	all.Mutex.Lock()
	defer all.Mutex.Unlock()
	gameSession, ok := all.Sessions[look_for]
	return gameSession, ok
}
func (allSessions *AllSession) IsSessionEmpty(playerId int) bool {
	allSessions.Mutex.Lock()
	defer allSessions.Mutex.Unlock()

	_, isFound := allSessions.Sessions[strconv.Itoa(playerId)]

	return !isFound
}

func (allSession *AllSession) Add(key string, gs *GameSession) {
	allSession.Mutex.Lock()
	defer allSession.Mutex.Unlock()

	allSession.Sessions[key] = gs
}
