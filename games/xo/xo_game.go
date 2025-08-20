package xo

import (
	"context"
	"sync"
	"time"

	"github.com/arian-nj/chibazi/database"
	gametype "github.com/arian-nj/chibazi/internals/game_type"
	"github.com/arian-nj/chibazi/internals/random"
	"github.com/arian-nj/chibazi/internals/socket"
	"github.com/arian-nj/chibazi/internals/utils"
	"gopkg.in/telebot.v4"
)

const MaxPlayerTime = time.Minute * 2

type XOGame struct { // of GameInterface type
	GameType gametype.GameType

	Board *XoBoard

	Players            []*XoPlayer
	CurrentPlayerIndex int

	Queries *database.Queries

	CancelGame context.CancelFunc
	Ctx        context.Context

	Actions     []Action
	Subscribers []XoSubscriber

	mu sync.Mutex
}

func NewXOGame(
	sessionCtx context.Context, gameType gametype.GameType,
	bot *telebot.Bot, queries *database.Queries) *XOGame {

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
		Board:    NewTicBoard(maxBoardSize, winSize),

		Queries: queries,

		Actions:     []Action{},
		Subscribers: []XoSubscriber{},
	}
}

func (game *XOGame) Register(subscriber XoSubscriber) {
	game.Subscribers = append(game.Subscribers, subscriber)
}

func (game *XOGame) PushAction(newAction Action) {
	game.Actions = append(game.Actions, newAction)
}

func (game *XOGame) InjectAction(newAction Action) {
	game.Actions = append([]Action{newAction}, game.Actions...)
}
func (game *XOGame) PopAction() Action {
	firstAction := game.Actions[0]
	game.Actions = game.Actions[1:]
	return firstAction
}

func (gs *XOGame) Apply(action Action) {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	action.Execute(gs)
	gs.Notify(action)
}

func (gs *XOGame) Notify(action Action) {
	for _, sub := range gs.Subscribers {
		sub.Update(gs, action) // pass both state + action
	}
}

func (g *XOGame) AddPlayer(name string, tgId int, socket *socket.Socket) {
	player := NewXoPlayer(name, tgId, socket)
	g.Players = append(g.Players, player)
}

func (cg *XOGame) GetCurrentPlayer() *XoPlayer {
	return cg.Players[cg.CurrentPlayerIndex]
}

func (cg *XOGame) IsPlayersTurn(senderId int) bool {
	return cg.GetCurrentPlayer().TgID == senderId
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

func (game *XOGame) StartGame() error {
	utils.RunBackgroundTask(func() {
		game.MonitorTimeout()
	})
	for _, player := range game.Players {
		now := time.Now()
		player.TurnStartedAt = now
	}
	game.Players[0].Move = X
	game.Players[1].Move = O

	startAction := NewStartAction()
	game.PushAction(startAction)
	// err = g.StartSocket()
	// if err != nil {
	// 	slog.Error("error starting game in web ", "error", err)
	// }

	return nil
}

func (game *XOGame) MonitorTimeout() {
	// ticker := time.NewTicker(time.Second * 1)
	// defer ticker.Stop()
	for {
		select {
		// case <-ticker.C:
		//
		// 	now := time.Now()
		// 	player := g.GetCurrentPlayer()
		// 	player.SpentTime += now.Sub(player.TurnStartedAt)
		// 	player.TurnStartedAt = time.Now()
		//
		// 	if player.SpentTime >= MaxPlayerTime {
		// 		g.NextPlayer()
		// 		err := g.TheEnd(bot, "\n برنده زمانی")
		// 		if err != nil {
		// 			slog.Error("error ending game with time out", "err", err)
		// 		}
		// 		return
		// 	}
		// 	if now.Sub(g.LastEdit) > time.Second*10 {
		// 		err := g.EditDuringGameBoard(bot)
		// 		if err != nil {
		// 			slog.Error("can't edit message in time monitor", "err", err)
		// 		}
		// 	}
		case <-game.Ctx.Done():
			return
		default:
			if len(game.Actions) > 0 {
				action := game.PopAction()
				game.Apply(action)
			}
		}
	}
}
