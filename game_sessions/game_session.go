package gamesessions

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/arian-nj/chibazi/database"
	gametype "github.com/arian-nj/chibazi/internals/game_type"
	humanplayer "github.com/arian-nj/chibazi/internals/human_player"
	"github.com/arian-nj/chibazi/internals/keybul"
	"github.com/arian-nj/chibazi/internals/socket"
	"github.com/arian-nj/chibazi/internals/utils"
	"gopkg.in/telebot.v4"
)

type SessionType string

const (
	PrivateSession SessionType = "private"
	RandomSession  SessionType = "random"
)

type SessionPlayer struct {
	ID     int
	TgID   int
	Name   string
	Socket *socket.Socket
}

func NewSessionPlayer(ID int, tgID int, name string) *SessionPlayer {
	return &SessionPlayer{
		ID:   ID,
		TgID: tgID,
		Name: name,
	}
}

type Game interface {
	Players() []*humanplayer.HumanPlayer
	AddPlayer(player *humanplayer.HumanPlayer)
	CallbackHandler(c telebot.Context) error
	GetContext() context.Context
	StartGame(bot telebot.API) error
	SendJoinPanelAddSender(telebot.Context) error
}

type GameSession struct {
	ID          int
	Bot         *telebot.Bot
	IsGameEnded bool
	Chat        *Chat
	Gametype    gametype.GameType
	GameState   Game

	Players []*SessionPlayer

	CreatedAt       time.Time
	ExpireDuaration time.Duration

	MsgChnl chan *SessionEvent

	Queries *database.Queries
}

func NewGameSession(allSession *AllSession, bot *telebot.Bot, Queries *database.Queries,
	gameType gametype.GameType, gameState Game,
	sessionId int) *GameSession {

	gs := &GameSession{
		ID:        sessionId,
		Bot:       bot,
		Gametype:  gameType,
		GameState: gameState,
		CreatedAt: time.Now(),
		Chat: &Chat{
			IsOn: true,
		},
		Players:         []*SessionPlayer{},
		ExpireDuaration: 2*time.Minute*2 + 30,
		MsgChnl:         make(chan *SessionEvent, 10),

		Queries: Queries,
	}

	utils.RunBackgroundTask(func() {
		gs.MonitorGameSession(allSession)
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
	utils.RunBackgroundTask(func() {
		_, err := gs.Queries.CreateSessionPlayer(context.Background(), database.CreateSessionPlayerParams{
			SessionID: gs.ID,
			PersonID:  player.ID,
		})
		if err != nil {
			slog.Error("can't add session player", "error", err)
		}
	})
}

func (gs *GameSession) MonitorGameSession(allSession *AllSession) {
	for {
		select {
		case newSEvent := <-gs.MsgChnl:
			switch newSEvent.Event.Type {
			case ChatType:
				gs.HandleWebChatMessage(newSEvent)
			}

		case <-gs.GameState.GetContext().Done():
			gs.IsGameEnded = true
			if gs.Chat.IsOn {
				expDur := 30 * time.Second

				text := fmt.Sprintf("چت تا %d ثانیه دیگه بسته میشه", int(expDur.Seconds()))
				for _, player := range gs.GameState.Players() {
					_, err := gs.Bot.Send(player, text)
					if err != nil {
						slog.Error("can't send end game chat message", "err", err)
					}
				}
				time.Sleep(expDur)
			}
			gs.CleanAndDisconnect(allSession)
			return
		}
	}
}

func (gs *GameSession) CleanAndDisconnect(allSession *AllSession) {
	if gs.Chat.IsOn {
		text := "چت قطع شد"
		for _, player := range gs.GameState.Players() {
			_, err := gs.Bot.Send(player, text, keybul.WelcomeReplyKeyboard)
			if err != nil {
				slog.Error("can't send chat ended message", "err", err)
			}
		}
	}

	allSession.Mutex.Lock()
	defer allSession.Mutex.Unlock()

	for _, player := range gs.GameState.Players() {
		delete(allSession.Sessions, strconv.Itoa(player.TgID))
	}
}
