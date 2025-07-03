package gamebot

import (
	"gopkg.in/telebot.v4"
)

type GameType string

const (
	TicTacToeGameType GameType = "ttt"
	DotBoxGameType    GameType = "dotbox"
)

type GameInterface interface {
	AddPlayer(player *HumanPlayer)
	GetPlayers() []*HumanPlayer
	SetMessageID(messageID string)

	GetGameType() GameType
	GetMaxPlayer() int
	StartGame(bot *telebot.Bot)
}

type Game struct {
	CurrentPlayerIndex int
	Players            []*HumanPlayer
	MessageId          string
}

func NewGame() *Game {
	return &Game{
		Players: []*HumanPlayer{},
	}
}
func (g *Game) SetMessageID(messageID string) {
	g.MessageId = messageID
}

func (g *Game) NextPlayer() {
	if g.CurrentPlayerIndex == len(g.Players)-1 {
		g.CurrentPlayerIndex = 0
		return
	}
	g.CurrentPlayerIndex += 1
}

func (g *Game) MessageSig() (string, int64) {
	return g.MessageId, 0
}

func (g *Game) AddPlayer(player *HumanPlayer) {
	g.Players = append(g.Players, player)
}
func (g *Game) GetPlayers() []*HumanPlayer {
	return g.Players
}
