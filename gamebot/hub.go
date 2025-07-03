package gamebot

import (
	"time"
)

type Lobby struct {
	Hubs HubsMap
}

func NewLobby() *Lobby {
	return &Lobby{
		Hubs: make(HubsMap),
	}
}

type HubsMap map[string]*Hub

func (hubsMap HubsMap) AddHub(hub *Hub) {
	hubsMap[hub.MessageId] = hub
}

type Hub struct {
	Game      GameInterface
	MessageId string
	CreatedAt time.Time
}

func NewHub(game GameInterface, messageId string) *Hub {
	return &Hub{
		MessageId: messageId,
		Game:      game,
		CreatedAt: time.Now(),
	}

}

func (hub *Hub) JoinGame(player *HumanPlayer, app *Application) bool {
	if len(hub.Game.GetPlayers()) >= 2 {
		return false
	}

	hub.Game.AddPlayer(player)
	if len(hub.Game.GetPlayers()) == hub.Game.GetMaxPlayer() {
		hub.Game.StartGame(app.Bot)
	}
	return true

}

type HumanPlayer struct {
	TgID int
	Name string
}

func NewHumanPlayer(tgID int, name string) *HumanPlayer {
	return &HumanPlayer{
		TgID: tgID,
		Name: name,
	}
}
