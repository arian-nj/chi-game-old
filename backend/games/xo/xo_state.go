package xo

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/arian-nj/chibazi/backend/database"
	"github.com/arian-nj/chibazi/backend/games/games"
	"github.com/arian-nj/chibazi/backend/internals/commander"
	"github.com/arian-nj/chibazi/backend/internals/random"
	"github.com/arian-nj/chibazi/backend/internals/socket"
	"github.com/arian-nj/chibazi/backend/internals/utils"
	"gopkg.in/telebot.v4"
)

const MAX_ALLOWED_TIME_INT = 60
const MAX_ALLOWED_TIME = MAX_ALLOWED_TIME_INT * time.Second

const (
	XOStartText = `❌ *دوز بازی* ⭕️`
	// ticRulesText = `
	// قوانین 🎮
	// یک سطر یا ستون یا قطر رو با علامتت پر کن`
)

func (game *XOState) RulesText() string {
	text := ""
	// text += "قوانین:\د"
	text += fmt.Sprintf("❕اندازه *%dX%d*\n", game.Board.MaxCellSize, game.Board.MaxCellSize)
	text += fmt.Sprintf("⚠️با یه خط *%d تایی* برنده ای", game.Board.WinSize)
	return text
}

type XOState struct { // of GameInterface type
	GameType games.GameType
	GameData *games.GameData

	Board *XoBoard

	Players            []*XoPlayer
	CurrentPlayerIndex int

	Queries *database.Queries

	CancelGame context.CancelFunc
	Ctx        context.Context

	*commander.Commander

	endCallback func()
}

func NewXOGame(
	sessionCtx context.Context, gameType games.GameType,
	queries *database.Queries) *XOState {

	maxBoardSize := 3
	winSize := 3
	if gameType == games.XOGameType5X5 {
		maxBoardSize = 5
		winSize = 4
	}
	randIndex := random.GenerateRandomNumber(2)

	ctx, cancel := context.WithCancel(sessionCtx)

	state := &XOState{
		CurrentPlayerIndex: randIndex,
		Players:            []*XoPlayer{},
		CancelGame:         cancel,
		Ctx:                ctx,

		GameType: gameType,
		Board:    NewTicBoard(maxBoardSize, winSize),

		Queries:   queries,
		Commander: commander.NewCommander(),
	}
	state.GameData = games.NewGameData(games.XOGameType3X3, XOStartText, state.RulesText(), 2)
	return state
}

func (gameState *XOState) GetGameData() *games.GameData {
	return gameState.GameData
}

func (gameState *XOState) OnEnd(endCallback func()) {
	gameState.endCallback = endCallback
}

func (gameState *XOState) monitorXoGame() {
	tickerDuration := time.Second * 1
	ticker := time.NewTicker(tickerDuration)
	lastSyncTime := time.Now().Add(time.Second * -10)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			now := time.Now()
			currentPlayer := gameState.CurrentPlayer()

			if currentPlayer.Timer.Spent() >= MAX_ALLOWED_TIME {
				newEndCommand := NewEndGameCommand(gameState, gameState.OpponentPlayer(), gameState.CurrentPlayer(), END_GAME_TIE)
				gameState.InjectCommand(newEndCommand)
			}
			if now.Sub(lastSyncTime) > time.Second*1 {
				newSyncCommand := NewSyncTimeCommand(gameState)
				gameState.PushCommand(newSyncCommand)
				lastSyncTime = now
			}
		case <-gameState.Ctx.Done():
			if gameState.endCallback != nil {
				gameState.endCallback()
			} else {
				slog.Error("there is not endCallback in xo game state")
			}
			return
		case <-gameState.CommandNotifire:
			if len(gameState.Commands) > 0 {
				com := gameState.PopCommand()
				gameState.ApplyCommand(com)
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
	g.CurrentPlayer().Timer.Stop()
	if g.CurrentPlayerIndex == len(g.Players)-1 {
		g.CurrentPlayerIndex = 0
	} else {
		g.CurrentPlayerIndex += 1
	}
	g.CurrentPlayer().Timer.Start()
}

// game interface
func (g *XOState) AddPlayer(id int, name string, tgId int, socket *socket.Socket) {
	player := NewXoPlayer(id, name, tgId, socket)
	g.Players = append(g.Players, player)
}

func (gameState *XOState) SubToTelegram(userID int, bot *telebot.Bot, ViaMessageId string) {
	foundPlayer := gameState.findByID(userID)
	if foundPlayer == nil && ViaMessageId == "" {
		slog.Error("no user found")
		return
	}
	gameState.Subscribe(newXOTelegramListener(foundPlayer, bot, ViaMessageId))
}

func (gameState *XOState) SubToSocket(ID int, newSocket *socket.Socket) func() {
	foundPlayer := gameState.findByID(ID)
	if foundPlayer == nil {
		return nil
	}

	foundPlayer.Socket = newSocket
	SocketSendGameState(gameState, foundPlayer)
	sListener := &SocketListener{
		Player: foundPlayer,
	}
	gameState.Subscribe(sListener)
	unregister := func() {
		if gameState != nil {
			gameState.Unsubscribe(sListener)
		}
	}
	return unregister
}

func (game *XOState) StartGame() error {

	utils.RunBackgroundTask(func() {
		game.monitorXoGame()
	})
	game.Players[0].Move = X
	game.Players[1].Move = O

	game.CurrentPlayer().Timer.Start()
	startAction := NewStartCommand(game)
	game.PushCommand(startAction)
	return nil
}
