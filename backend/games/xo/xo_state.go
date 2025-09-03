package xo

import (
	"context"
	"sync"
	"time"

	"github.com/arian-nj/chibazi/backend/database"
	gametype "github.com/arian-nj/chibazi/backend/internals/game_type"
	"github.com/arian-nj/chibazi/backend/internals/random"
	"github.com/arian-nj/chibazi/backend/internals/socket"
	"github.com/arian-nj/chibazi/backend/internals/utils"
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

	Commands     []Command
	DoneCommands []Command
	CommandCh    chan any

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

		Commands:     []Command{},
		DoneCommands: []Command{},
		CommandCh:    make(chan any, 6),

		Subscribers: []XoSubscriber{},
	}
}

func (game *XOGame) monitorTimeout() {
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
		case <-game.CommandCh:
			if len(game.Commands) > 0 {
				action := game.popCommand()
				game.applyCommand(action)
			} else {
				time.Sleep(50 * time.Millisecond)
			}
		}
	}
}

// helper functions
func (game *XOGame) findByTelegramID(telegramID int) *XoPlayer {
	for _, p := range game.Players {
		if telegramID == p.TelegramID {
			return p
		}
	}
	return nil
}

func (game *XOGame) findByID(telegramID int) *XoPlayer {
	for _, p := range game.Players {
		if telegramID == p.ID {
			return p
		}
	}
	return nil
}

func (cg *XOGame) getCurrentPlayer() *XoPlayer {
	return cg.Players[cg.CurrentPlayerIndex]
}

func (g *XOGame) nextPlayer() {
	if g.CurrentPlayerIndex == len(g.Players)-1 {
		g.CurrentPlayerIndex = 0
	} else {
		g.CurrentPlayerIndex += 1
	}
	g.getCurrentPlayer().TurnStartedAt = time.Now()
}

// Commands
func (game *XOGame) pushCommand(newCommand Command) {
	game.Commands = append(game.Commands, newCommand)
	game.CommandCh <- nil
}

func (game *XOGame) injectCommand(newAction Command) {
	game.Commands = append([]Command{newAction}, game.Commands...)
	game.CommandCh <- nil
}
func (game *XOGame) popCommand() Command {
	firstAction := game.Commands[0]
	game.Commands = game.Commands[1:]
	return firstAction
}

func (gs *XOGame) applyCommand(newCommand Command) {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	newCommand.Execute(gs)
	gs.Notify(newCommand)
	gs.DoneCommands = append(gs.DoneCommands, newCommand)
}

// Subscription
func (gs *XOGame) Notify(command Command) {
	for _, sub := range gs.Subscribers {
		sub.Update(gs, command) // pass both state + action
	}
}

func (game *XOGame) Register(subscriber XoSubscriber) {
	game.Subscribers = append(game.Subscribers, subscriber)
}

// game interface
func (g *XOGame) AddPlayer(id int, name string, tgId int, socket *socket.Socket) {
	player := NewXoPlayer(id, name, tgId, socket)
	g.Players = append(g.Players, player)
}

func (cg *XOGame) GetContext() context.Context {
	return cg.Ctx
}

func (game *XOGame) SetPlayerSocket(ID int, newSocket *socket.Socket) {
	foundPlayer := game.findByID(ID)
	if foundPlayer != nil {
		foundPlayer.Socket = newSocket
	}
}

func (game *XOGame) StartGame() error {
	utils.RunBackgroundTask(func() {
		game.monitorTimeout()
	})
	for _, player := range game.Players {
		now := time.Now()
		player.TurnStartedAt = now
	}
	game.Players[0].Move = X
	game.Players[1].Move = O

	startAction := NewStartCommand()
	game.pushCommand(startAction)
	// err = g.StartSocket()
	// if err != nil {
	// 	slog.Error("error starting game in web ", "error", err)
	// }

	return nil
}
