package conn4

import (
	"context"
	"log/slog"
	"time"

	"github.com/arian-nj/chibazi/backend/database"
	conn4_core "github.com/arian-nj/chibazi/backend/games/conn4/core"
	"github.com/arian-nj/chibazi/backend/games/games"
	"github.com/arian-nj/chibazi/backend/internals/commander"
	"github.com/arian-nj/chibazi/backend/internals/random"
	"github.com/arian-nj/chibazi/backend/internals/socket"
	"github.com/arian-nj/chibazi/backend/internals/utils"
	"gopkg.in/telebot.v4"
)

const MAX_ALLOWED_TIME_INT = 10
const MAX_ALLOWED_TIME = MAX_ALLOWED_TIME_INT * time.Second

const (
	Conn4StartText = `❌ *دوز بازی* ⭕️`
	Conn4RulesText = `قوانین`
	// قوانین 🎮
	// یک سطر یا ستون یا قطر رو با علامتت پر کن`
)

type Conn4State struct { // of GameInterface type
	GameType games.GameType
	GameData *games.GameData

	Board *conn4_core.Conn4Board

	Players            []*Conn4Player
	CurrentPlayerIndex int

	Queries *database.Queries

	CancelGame context.CancelFunc
	Ctx        context.Context

	*commander.Commander

	endCallback func()
}

func NewConn4State(
	sessionCtx context.Context, gameType games.GameType,
	queries *database.Queries) *Conn4State {

	randIndex := random.GenerateRandomNumber(2)

	ctx, cancel := context.WithCancel(sessionCtx)

	state := &Conn4State{
		CurrentPlayerIndex: randIndex,
		Players:            []*Conn4Player{},
		CancelGame:         cancel,
		Ctx:                ctx,

		GameType: gameType,
		Board:    conn4_core.NewConn4Board(),

		Queries:   queries,
		Commander: commander.NewCommander(),
	}
	state.GameData = games.NewGameData(games.Conn4GameType, Conn4StartText,
		"قوانین",
		2,
	)
	return state
}

func (gameState *Conn4State) GetGameData() *games.GameData {
	return gameState.GameData
}

func (gameState *Conn4State) OnEnd(endCallback func()) {
	gameState.endCallback = endCallback
}
func (gameState *Conn4State) monitorXoGame() {
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
				slog.Info("add timeout end end command")
				newEndCommand := NewEndGameCommand(gameState, gameState.OpponentPlayer(), gameState.CurrentPlayer(), END_GAME_TIMEOUT, []int{})
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
func (game *Conn4State) findByTelegramID(telegramID int) *Conn4Player {
	for _, p := range game.Players {
		if telegramID == p.TelegramID {
			return p
		}
	}
	return nil
}

func (game *Conn4State) findByID(telegramID int) *Conn4Player {
	for _, p := range game.Players {
		if telegramID == p.ID {
			return p
		}
	}
	return nil
}

func (game *Conn4State) CurrentPlayer() *Conn4Player {
	return game.Players[game.CurrentPlayerIndex]
}

func (gameState *Conn4State) OpponentPlayer() *Conn4Player {
	index := 0
	if gameState.CurrentPlayerIndex == 0 {
		index = 1
	}
	return gameState.Players[index]
}

func (g *Conn4State) nextPlayer() {
	g.CurrentPlayer().Timer.Stop()
	if g.CurrentPlayerIndex == len(g.Players)-1 {
		g.CurrentPlayerIndex = 0
	} else {
		g.CurrentPlayerIndex += 1
	}
	g.CurrentPlayer().Timer.Start()
}

func (g *Conn4State) AddPlayer(id int, name string, tgId int, socket *socket.Socket) {
	player := NewConn4Player(id, name, tgId, socket)
	g.Players = append(g.Players, player)
}

func (gameState *Conn4State) SubToTelegram(userID int, bot *telebot.Bot, ViaMessageId string) {
	foundPlayer := gameState.findByID(userID)
	if foundPlayer == nil && ViaMessageId == "" {
		slog.Error("no user found")
		return
	}
	gameState.Subscribe(newConn4TelegramListener(foundPlayer, bot, ViaMessageId))
}

func (gameState *Conn4State) SubToSocket(ID int, newSocket *socket.Socket) func() {
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

func (game *Conn4State) StartGame() error {
	utils.RunBackgroundTask(func() {
		game.monitorXoGame()
	})
	game.Players[0].Move = conn4_core.One
	game.Players[1].Move = conn4_core.Two

	game.CurrentPlayer().Timer.Start()
	startAction := NewStartCommand(game)
	game.PushCommand(startAction)
	return nil
}
