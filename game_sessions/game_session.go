package gamesessions

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/arian-nj/chibazi/database"
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
		Name: keybul.EscapeReserved(name),
	}
}

type GameSession struct {
	ID int

	Bot     *telebot.Bot
	Queries *database.Queries

	IsGameEnded bool
	Chat        *Chat
	GameState   Game

	MsgChnl chan *SessionEvent

	Players []*SessionPlayer

	CancelSession context.CancelFunc
	SessionCtx    context.Context

	CreatedAt       time.Time
	ExpireDuaration time.Duration
}

func NewGameSession(bot *telebot.Bot, Queries *database.Queries, sessionId int) *GameSession {

	ctx, cancel := context.WithCancel(context.Background())
	gs := &GameSession{
		ID:        sessionId,
		Bot:       bot,
		CreatedAt: time.Now(),
		Chat: &Chat{
			IsOn: true,
		},
		Players:         []*SessionPlayer{},
		ExpireDuaration: 2*time.Minute*2 + 30,
		MsgChnl:         make(chan *SessionEvent, 10),

		Queries: Queries,

		CancelSession: cancel,
		SessionCtx:    ctx,
	}
	return gs
}

func (gs *GameSession) RunBgTask(allSession *AllSession) {
	utils.RunBackgroundTask(func() {
		gs.MonitorGameSession(allSession)
	})
}

func (gs *GameSession) StartGame() error {
	if gs.GameState == nil {
		return errors.New("failed to start game field GameState is nil")
	}

	for _, player := range gs.Players {
		gs.GameState.AddPlayer(player.Name, player.TgID, player.Socket)
	}

	return gs.GameState.StartGame()
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
			case socket.ChatEventType:
				gs.HandleWebChatMessage(newSEvent)
			case socket.GameEventType:
				// gs.SocketRouter()
			}

		case <-gs.GameState.GetContext().Done():
			gs.IsGameEnded = true
			if gs.Chat.IsOn {
				expDur := 30 * time.Second

				text := fmt.Sprintf("چت تا %d ثانیه دیگه بسته میشه", int(expDur.Seconds()))
				for _, player := range gs.Players {
					_, err := gs.Bot.Send(&telebot.User{ID: int64(player.TgID)}, text)
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
		for _, player := range gs.Players {
			_, err := gs.Bot.Send(&telebot.User{ID: int64(player.TgID)}, text, keybul.WelcomeReplyKeyboard)
			if err != nil {
				slog.Error("can't send chat ended message", "err", err)
			}
		}
	}

	allSession.Mutex.Lock()
	defer allSession.Mutex.Unlock()

	for _, player := range gs.Players {
		delete(allSession.Sessions, strconv.Itoa(player.TgID))
	}
}
