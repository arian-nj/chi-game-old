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

const MaxAllowedTimeSecond = 120

type XOState struct { // of GameInterface type
	GameType gametype.GameType

	Board *XoBoard

	Players            []*XoPlayer
	CurrentPlayerIndex int

	Queries *database.Queries

	CancelGame context.CancelFunc
	Ctx        context.Context

	Commands        []Command
	DoneCommands    []Command
	CommandNotifire chan any

	Subscribers []XoSubscriber

	mu sync.Mutex
}

func NewXOGame(
	sessionCtx context.Context, gameType gametype.GameType,
	bot *telebot.Bot, queries *database.Queries) *XOState {

	maxBoardSize := 3
	winSize := 3
	if gameType == gametype.XOGameType5X5 {
		maxBoardSize = 5
		winSize = 4
	}
	randIndex := random.GenerateRandomNumber(2)

	ctx, cancel := context.WithCancel(sessionCtx)

	return &XOState{

		CurrentPlayerIndex: randIndex,
		Players:            []*XoPlayer{},
		CancelGame:         cancel,
		Ctx:                ctx,

		GameType: gameType,
		Board:    NewTicBoard(maxBoardSize, winSize),

		Queries: queries,

		Commands:        []Command{},
		DoneCommands:    []Command{},
		CommandNotifire: make(chan any, 6),

		Subscribers: []XoSubscriber{},
	}
}

func (gameState *XOState) monitorXoGame() {
	tickerDuration := time.Millisecond * 500
	ticker := time.NewTicker(tickerDuration)
	lastSyncTime := time.Now()
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			now := time.Now()
			currentPlayer := gameState.CurrentPlayer()
			currentPlayer.SpentTime += now.Sub(currentPlayer.LastTurnStartedAt)
			currentPlayer.LastTurnStartedAt = now

			if currentPlayer.SpentTime >= MaxAllowedTimeSecond*time.Second {
				newEndCommand := NewEndGameCommand(gameState.OpponentPlayer(), "\n برنده زمانی")
				gameState.injectCommand(newEndCommand)
				return
			}
			if now.Sub(lastSyncTime) > time.Second*5 {
				newSyncCommand := NewSyncTimeCommand()
				gameState.pushCommand(newSyncCommand)
			}
		case <-gameState.Ctx.Done():
			return
		case <-gameState.CommandNotifire:
			if len(gameState.Commands) > 0 {
				action := gameState.popCommand()
				gameState.applyCommand(action)
			} else {
				time.Sleep(50 * time.Millisecond)
			}
		}
	}
}

// helper functions
func (game *XOState) findByTelegramID(telegramID int) *XoPlayer {
	for _, p := range game.Players {
		if telegramID == p.TelegramID {
			return p
		}
	}
	return nil
}

func (game *XOState) findByID(telegramID int) *XoPlayer {
	for _, p := range game.Players {
		if telegramID == p.ID {
			return p
		}
	}
	return nil
}

func (cg *XOState) CurrentPlayer() *XoPlayer {
	return cg.Players[cg.CurrentPlayerIndex]
}

func (gameState *XOState) OpponentPlayer() *XoPlayer {
	index := 0
	if gameState.CurrentPlayerIndex == 0 {
		index = 1
	}
	return gameState.Players[index]
}

func (g *XOState) nextPlayer() {
	if g.CurrentPlayerIndex == len(g.Players)-1 {
		g.CurrentPlayerIndex = 0
	} else {
		g.CurrentPlayerIndex += 1
	}
	g.CurrentPlayer().LastTurnStartedAt = time.Now()
}

// Commands
func (game *XOState) pushCommand(newCommand Command) {
	game.Commands = append(game.Commands, newCommand)
	game.CommandNotifire <- nil
}

func (game *XOState) injectCommand(newAction Command) {
	game.Commands = append([]Command{newAction}, game.Commands...)
	game.CommandNotifire <- nil
}
func (game *XOState) popCommand() Command {
	firstAction := game.Commands[0]
	game.Commands = game.Commands[1:]
	return firstAction
}

func (gs *XOState) applyCommand(newCommand Command) {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	newCommand.Execute(gs)
	gs.Notify(newCommand)
	gs.DoneCommands = append(gs.DoneCommands, newCommand)
}

// Subscription
func (gs *XOState) Notify(command Command) {
	for _, sub := range gs.Subscribers {
		sub.Update(gs, command) // pass both state + action
	}
}

func (game *XOState) Register(subscriber XoSubscriber) {
	game.Subscribers = append(game.Subscribers, subscriber)
}

// game interface
func (g *XOState) AddPlayer(id int, name string, tgId int, socket *socket.Socket) {
	player := NewXoPlayer(id, name, tgId, socket)
	g.Players = append(g.Players, player)
}

func (cg *XOState) GetContext() context.Context {
	return cg.Ctx
}

func (game *XOState) SetPlayerSocket(ID int, newSocket *socket.Socket) {
	foundPlayer := game.findByID(ID)
	if foundPlayer != nil {
		foundPlayer.Socket = newSocket
	}
}

func (game *XOState) StartGame() error {
	utils.RunBackgroundTask(func() {
		game.monitorXoGame()
	})
	for _, player := range game.Players {
		now := time.Now()
		player.LastTurnStartedAt = now
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
