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

// func (hp HubPlayers) ChooseRandom() (*HumanPlayer, error) {
// 	if len(hp) <= 0 {
// 		return nil, fmt.Errorf("erro there is 0 players")
// 	}
// 	randomIndex := random.GenerateRandomNumber(len(hp))
// 	return hp[randomIndex], nil
// }
//

type Hub struct {
	Name      string
	Game      GameInterface
	MessageId string
	Players   HubPlayers
	CreatedAt time.Time
}

func (hub *Hub) MessageSig() (string, int64) {
	return hub.MessageId, 0
}

func NewHub(name string, game GameInterface, messageId string) *Hub {
	return &Hub{
		Name:      name,
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
	TicTacToe GameType = "tictactoe"
)

type GameInterface interface {
	GetGameType() GameType
	StartGame(app *Application)
	JoinGame(app *Application, player *HumanPlayer) bool
	EndGame()
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
