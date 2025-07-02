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

type GameType string

const (
	TicTacToe GameType = "ttt"
)

type GameInterface interface {
	GetGameType() GameType
	StartGame(app *Application)
	JoinGame(app *Application, player *HumanPlayer) bool
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
