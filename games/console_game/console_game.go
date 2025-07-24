package consolegame

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	gametype "github.com/arian-nj/chibazi/internals/game_type"
	humanplayer "github.com/arian-nj/chibazi/internals/human_player"
	"github.com/arian-nj/chibazi/internals/keybul"
	"gopkg.in/telebot.v4"
)

type ConsoleGame struct {
	GameType   gametype.GameType
	CancelGame context.CancelFunc
	Ctx        context.Context

	ViaMessageId string // Via Bots

	players []*humanplayer.HumanPlayer

	LastEdit           time.Time
	CurrentPlayerIndex int
}

func NewConsoleGame(gameType gametype.GameType, currentPlayerIndex int) *ConsoleGame {
	ctx, cancel := context.WithCancel(context.Background())
	return &ConsoleGame{
		GameType:           gameType,
		CancelGame:         cancel,
		Ctx:                ctx,
		CurrentPlayerIndex: currentPlayerIndex,

		players: []*humanplayer.HumanPlayer{},
	}
}

func (cg *ConsoleGame) MessageSig() (string, int64) {
	return cg.ViaMessageId, 0
}

func (cg *ConsoleGame) GetContext() context.Context {
	return cg.Ctx
}

func (cg *ConsoleGame) GetCurrentPlayer() *humanplayer.HumanPlayer {
	return cg.players[cg.CurrentPlayerIndex]
}
func (cg *ConsoleGame) Players() []*humanplayer.HumanPlayer {
	return cg.players
}
func (g *ConsoleGame) AddPlayer(player *humanplayer.HumanPlayer) {
	g.players = append(g.players, player)
}

func (g *ConsoleGame) NextPlayer() {
	if g.CurrentPlayerIndex == len(g.Players())-1 {
		g.CurrentPlayerIndex = 0
	} else {
		g.CurrentPlayerIndex += 1
	}
	g.GetCurrentPlayer().TurnStartedAt = time.Now()

}

func (g *ConsoleGame) Edit(bot telebot.API, msg telebot.Editable, text string, keyboard *telebot.ReplyMarkup) error {
	g.LastEdit = time.Now()
	if g.ViaMessageId != "" {
		err := keybul.EditGameMessage(bot, g, text, keyboard)
		if err != nil {
			return fmt.Errorf("can't edit via message %w", err)
		}
		return nil
	} else {
		for _, p := range g.Players() {
			if p.MessageID == 0 {
				msg, err := bot.Send(p, "game")
				if err != nil {
					slog.Error("can't send player message ", "error", err)
				}
				p.SetMessageSig(msg.ID)
			}
			err := keybul.EditGameMessage(bot, p, text, keyboard)
			if err != nil {
				slog.Error("can't edit player message ", "error", err)
			}
		}

	}
	return nil
}
