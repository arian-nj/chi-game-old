package xoconsole

import (
	"context"
	"log/slog"
	"time"

	"github.com/arian-nj/chibazi/database"
	gametype "github.com/arian-nj/chibazi/internals/game_type"
	"github.com/arian-nj/chibazi/internals/keybul"
	"github.com/arian-nj/chibazi/internals/random"
	"github.com/arian-nj/chibazi/internals/socket"
	"github.com/arian-nj/chibazi/internals/xo_core"
	"gopkg.in/telebot.v4"
)

const MaxPlayerTime = time.Minute * 2

type XOGame struct { // of GameInterface type
	GameType gametype.GameType
	Bot      *telebot.Bot

	Board *xo_core.XoBoard

	Players            []*XoPlayer
	CurrentPlayerIndex int

	ViaMessageId string // Via Bots
	LastEdit     time.Time

	Queries *database.Queries

	CancelGame context.CancelFunc
	Ctx        context.Context
}

func NewXOGame(sessionCtx context.Context, gameType gametype.GameType, bot *telebot.Bot, queries *database.Queries) *XOGame {
	maxBoardSize := 3
	winSize := 3
	if gameType == gametype.XOGameType5X5 {
		maxBoardSize = 5
		winSize = 4
	}
	randIndex := random.GenerateRandomNumber(2)

	ctx, cancel := context.WithCancel(sessionCtx)
	return &XOGame{

		CurrentPlayerIndex: randIndex,
		Players:            []*XoPlayer{},
		CancelGame:         cancel,
		Ctx:                ctx,

		GameType: gameType,
		Board:    xo_core.NewTicBoard(maxBoardSize, winSize),

		Queries: queries,
		Bot:     bot,
	}
}
func (g *XOGame) AddPlayer(name string, tgId int, socket *socket.Socket) {
	player := NewXoPlayer(name, tgId, socket)
	g.Players = append(g.Players, player)
}

func (cg *XOGame) GetCurrentPlayer() *XoPlayer {
	return cg.Players[cg.CurrentPlayerIndex]
}
func (cg *XOGame) GetOpponentPlayer() *XoPlayer {
	if cg.CurrentPlayerIndex == 0 {
		return cg.Players[1]
	}
	return cg.Players[0]
}

func (g *XOGame) NextPlayer() {
	if g.CurrentPlayerIndex == len(g.Players)-1 {
		g.CurrentPlayerIndex = 0
	} else {
		g.CurrentPlayerIndex += 1
	}
	g.GetCurrentPlayer().TurnStartedAt = time.Now()

}

func (cg *XOGame) GetContext() context.Context {
	return cg.Ctx
}

func (g *XOGame) SetPlayerSocket(tgId int, socket *socket.Socket) {
	for _, player := range g.Players {
		if player.TgID == tgId {
			player.Socket = socket
			return
		}
	}

}

func (g *XOGame) StartGame() error {
	err := g.StartGameBot()
	if err != nil {
		slog.Error("error starting game in bot ", "error", err)
	}

	err = g.StartSocket()
	if err != nil {
		slog.Error("error starting game in web ", "error", err)
	}

	for _, player := range g.Players {
		now := time.Now()
		player.TurnStartedAt = now
	}
	g.Players[0].Move = xo_core.X
	g.Players[1].Move = xo_core.O
	return nil
}

func (g *XOGame) TheEnd(bot telebot.API, additionalText string) error {
	g.CancelGame()
	text := g.EndGameText() + g.WinGameText() + additionalText
	err := g.Edit(bot, g, text,
		keybul.CreateInlineKeyboard(
			keybul.CreateBotNameInlineButton(),
			keybul.EndGameInlineKeyboard(g.ViaMessageId != ""),
		),
	)
	return err
}

func (g *XOGame) TieGame(bot telebot.API) error {
	text := g.EndGameText() + "\nبازی مساوی شد"

	err := g.Edit(bot, g, text,
		keybul.CreateInlineKeyboard(
			keybul.CreateBotNameInlineButton(),
			keybul.EndGameInlineKeyboard(g.ViaMessageId != ""),
		),
	)
	return err
}

// func (g *XOGame) MonitorTimeout(bot telebot.API) {
// 	ticker := time.NewTicker(time.Second * 1)
// 	defer ticker.Stop()
// 	for {
// 		select {
// 		case <-ticker.C:
//
// 			now := time.Now()
// 			player := g.GetCurrentPlayer()
// 			player.SpentTime += now.Sub(player.TurnStartedAt)
// 			player.TurnStartedAt = time.Now()
//
// 			if player.SpentTime >= MaxPlayerTime {
// 				g.NextPlayer()
// 				err := g.TheEnd(bot, "\n برنده زمانی")
// 				if err != nil {
// 					slog.Error("error ending game with time out", "err", err)
// 				}
// 				return
// 			}
// 			if now.Sub(g.LastEdit) > time.Second*10 {
// 				err := g.EditDuringGameBoard(bot)
// 				if err != nil {
// 					slog.Error("can't edit message in time monitor", "err", err)
// 				}
// 			}
// 		case <-g.Ctx.Done():
// 			return
// 		}
// 	}
// }
