package gamebot

import "time"

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

type HubPlayers []*HumanPlayer

type Hub struct {
	Game      GameInterface
	MessageId string
	Players   HubPlayers
	CreatedAt time.Time
}

func (hub *Hub) MessageSig() (string, int64) {
	return hub.MessageId, 0
}

func NewHub(game GameInterface, messageId string) *Hub {
	return &Hub{
		MessageId: messageId,
		Players:   make([]*HumanPlayer, 0),
		Game:      game,
		CreatedAt: time.Now(),
	}

}

func (hub *Hub) AddPlayer(player *HumanPlayer) {
	hub.Players = append(hub.Players, player)
}

func (hub *Hub) JoinGame(player *HumanPlayer, app *Application) bool {
	if len(hub.Players) >= 2 {
		return false
	}
	hub.AddPlayer(player)
	if len(hub.Players) == 2 {
		hub.Game.StartGame(app)
	}
	return true

}

type GameType string

const (
	TicTacToeGameType GameType = "ttt"
	DotBoxGameType    GameType = "dotbox"
)

type GameInterface interface {
	GetGameType() GameType
	GetMaxPlayer() int
	StartGame(app *Application)
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
