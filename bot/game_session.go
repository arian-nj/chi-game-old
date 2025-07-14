package bot

import (
	"fmt"
	"sync"
	"time"

	consoleplayer "github.com/arian-nj/chibazi/internals/console_player"
	gametype "github.com/arian-nj/chibazi/internals/game_type"
	"gopkg.in/telebot.v4"
)

type Game interface {
	Players() []*consoleplayer.ConsolePlayer
	CallbackHandler(c telebot.Context) error
}

type AllSession struct {
	Sessions map[string]*GameSession
	Mutex    sync.Mutex
}

type GameSession struct {
	Gametype  gametype.GameType
	ChatState bool

	GameState Game

	CreatedAt time.Time
}

func NewGameSession(gameType gametype.GameType, gameState Game) *GameSession {
	return &GameSession{
		Gametype:  gameType,
		GameState: gameState,
		ChatState: true,
		CreatedAt: time.Now(),
	}
}

func (g *GameSession) HandleChatMessage(bot telebot.API, senderID int, text string) error {
	players := g.GameState.Players()
	var senderPlayer *consoleplayer.ConsolePlayer
	var recieverPlayer *consoleplayer.ConsolePlayer

	for _, p := range players {
		if p.TgID != senderID {
			recieverPlayer = p
		} else {
			senderPlayer = p
		}
	}

	_, err := bot.Send(&telebot.User{ID: int64(recieverPlayer.TgID)},
		fmt.Sprintf("*_%s:_* %s", senderPlayer.Name, text), telebot.ModeMarkdownV2)
	return err
}
