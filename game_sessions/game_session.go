package gamesessions

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"sync"
	"time"

	gametype "github.com/arian-nj/chibazi/internals/game_type"
	humanplayer "github.com/arian-nj/chibazi/internals/human_player"
	"github.com/arian-nj/chibazi/internals/utils"
	"gopkg.in/telebot.v4"
)

type Game interface {
	Players() []*humanplayer.HumanPlayer
	CallbackHandler(c telebot.Context) error
	GetContext() context.Context
}

type AllSession struct {
	Sessions map[string]*GameSession
	Mutex    sync.Mutex
}

func (allSessions *AllSession) IsSessionEmpty(playerId int) bool {
	_, isFound := allSessions.Sessions[strconv.Itoa(playerId)]
	return !isFound
}

func (allSession *AllSession) Add(key string, gs *GameSession) {
	allSession.Mutex.Lock()
	defer allSession.Mutex.Unlock()

	allSession.Sessions[key] = gs
}

type GameSession struct {
	Bot         *telebot.Bot
	IsGameEnded bool
	Chat        *Chat
	Gametype    gametype.GameType
	GameState   Game

	CreatedAt       time.Time
	ExpireDuaration time.Duration
}

func NewGameSession(allSession *AllSession, bot *telebot.Bot, gameType gametype.GameType, gameState Game) *GameSession {
	gs := &GameSession{
		Bot:       bot,
		Gametype:  gameType,
		GameState: gameState,
		CreatedAt: time.Now(),
		Chat: &Chat{
			IsChatOn: true,
		},
		ExpireDuaration: 2*time.Minute*2 + 30,
	}

	utils.RunBackgroundTask(func() {
		gs.MonitorGame(allSession)
	})

	return gs
}
func (gs *GameSession) MonitorGame(allSession *AllSession) {
	<-gs.GameState.GetContext().Done()
	gs.IsGameEnded = true
	if gs.Chat.IsChatOn {
		expDur := 30 * time.Second
		gs.ExpireDuaration = time.Since(gs.CreatedAt) + expDur

		text := fmt.Sprintf("چت تا %d ثانیه دیگه بسته میشه", int(expDur.Seconds()))
		for _, player := range gs.GameState.Players() {
			_, err := gs.Bot.Send(player, text)
			if err != nil {
				slog.Error("can't send end game chat message", "err", err)
			}
		}
		time.Sleep(gs.ExpireDuaration)
	}
	gs.CleanAndDisconnect(allSession)
}

func (gs *GameSession) CleanAndDisconnect(allSession *AllSession) {
	// if gs.IsChatOn {
	// 	text := "چت قطع شد"
	// 	for _, player := range gs.GameState.Players() {
	// 		_, err := gs.Bot.Send(player, text, welcomeReplyKeyboard)
	// 		if err != nil {
	// 			slog.Error("can't send chat ended message", "err", err)
	// 		}
	// 	}
	// }
	//
	// allSession.Mutex.Lock()
	// defer allSession.Mutex.Unlock()
	//
	// for _, player := range gs.GameState.Players() {
	// 	delete(allSession.Sessions, strconv.Itoa(player.TgID))
	// }
}
