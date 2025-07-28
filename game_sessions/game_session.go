package gamesessions

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	gametype "github.com/arian-nj/chibazi/internals/game_type"
	humanplayer "github.com/arian-nj/chibazi/internals/human_player"
	"github.com/arian-nj/chibazi/internals/socket"
	"github.com/arian-nj/chibazi/internals/utils"
	"gopkg.in/telebot.v4"
)

type Game interface {
	Players() []*humanplayer.HumanPlayer
	AddPlayer(player *humanplayer.HumanPlayer)
	CallbackHandler(c telebot.Context) error
	GetContext() context.Context
	StartGame(bot telebot.API) error
}

type SessionPlayer struct {
	TgID   int
	Name   string
	Socket *socket.Socket
}

func NewSessionPlayer(tgID int, name string) *SessionPlayer {
	return &SessionPlayer{
		TgID: tgID,
		Name: name,
	}
}

type GameSession struct {
	Bot         *telebot.Bot
	IsGameEnded bool
	Chat        *Chat
	Gametype    gametype.GameType
	GameState   Game

	Players []*SessionPlayer

	CreatedAt       time.Time
	ExpireDuaration time.Duration

	MsgChnl chan *SessionEvent
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
		Players:         []*SessionPlayer{},
		ExpireDuaration: 2*time.Minute*2 + 30,
	}

	utils.RunBackgroundTask(func() {
		gs.MonitorGame(allSession)
	})

	return gs
}

func (gs *GameSession) StartGame() error {
	if gs.GameState == nil {
		return errors.New("failed to start game field GameState is nil")
	}

	for _, player := range gs.Players {
		newPlayer := humanplayer.NewHumanPlayer(player.Name, player.TgID)
		gs.GameState.AddPlayer(newPlayer)
	}

	return gs.GameState.StartGame(gs.Bot)
}

func (gs *GameSession) AddPlayer(player *SessionPlayer) {
	gs.Players = append(gs.Players, player)
}

func (gs *GameSession) MonitorGame(allSession *AllSession) {
	select {
	case newSEvent := <-gs.MsgChnl:
		switch newSEvent.Event.Type {
		case ChatType:
			gs.HandleWebChatMessage(newSEvent)
		}
	case <-gs.GameState.GetContext().Done():
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
