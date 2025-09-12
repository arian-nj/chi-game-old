package gamesessions

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/arian-nj/chibazi/backend/database"
	"github.com/arian-nj/chibazi/backend/games/game"
	sessionv1 "github.com/arian-nj/chibazi/backend/gen/session/v1"
	"github.com/arian-nj/chibazi/backend/internals/commander"
	"github.com/arian-nj/chibazi/backend/internals/keybul"
	"github.com/arian-nj/chibazi/backend/internals/socket"
	"github.com/arian-nj/chibazi/backend/internals/utils"
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

type Chat struct {
	IsOn bool
}
type GameSession struct {
	ID int

	Bot     *telebot.Bot
	Queries *database.Queries

	IsGameEnded bool
	Chat        Chat
	GameState   game.Game

	MsgChnl chan *SessionEvent

	Players []*SessionPlayer

	CancelSession context.CancelFunc
	SessionCtx    context.Context

	CreatedAt       time.Time
	ExpireDuaration time.Duration

	*commander.Commander
}

func NewGameSession(bot *telebot.Bot, Queries *database.Queries, sessionId int) *GameSession {

	ctx, cancel := context.WithCancel(context.Background())
	gs := &GameSession{
		ID:        sessionId,
		Bot:       bot,
		CreatedAt: time.Now(),
		Chat: Chat{
			IsOn: true,
		},
		Players:         []*SessionPlayer{},
		ExpireDuaration: 2*time.Minute*2 + 30,
		MsgChnl:         make(chan *SessionEvent, 10),

		Queries: Queries,

		CancelSession: cancel,
		SessionCtx:    ctx,

		Commander: commander.NewCommander(),
	}
	return gs
}

func (session *GameSession) RunBgTask(allSession *AllSession) {
	utils.RunBackgroundTask(func() {
		session.MonitorGameSession(allSession)
	})
}

func (session *GameSession) StartGame() error {
	if session.GameState == nil {
		return errors.New("failed to start game field GameState is nil")
	}

	for _, player := range session.Players {
		session.GameState.AddPlayer(player.ID, player.Name, player.TgID, player.Socket)
	}

	for _, suber := range session.Subscribers {
		if listener, ok := suber.(*SessionTelegramListener); ok {
			session.GameState.SubToTelegram(listener.Bot, listener.ViaMessageId)
		}
	}
	session.GameState.SubToSocket()

	return session.GameState.StartGame()
}

func (session *GameSession) AddSessionPlayer(player *SessionPlayer) {
	session.Players = append(session.Players, player)
	utils.RunBackgroundTask(func() {
		_, err := session.Queries.CreateSessionPlayer(context.Background(), database.CreateSessionPlayerParams{
			SessionID: session.ID,
			PersonID:  player.ID,
		})
		if err != nil {
			slog.Error("can't add session player", "error", err)
		}
	})
}

func (session *GameSession) MonitorGameSession(allSession *AllSession) {
	for {
		select {
		case newSessionEvent := <-session.MsgChnl:
			switch newSessionEvent.Event.Content.(type) {
			case *sessionv1.SessionMessage_ChatReq:
				session.SocketRequestSendMsg(newSessionEvent.Player, newSessionEvent.Event.GetChatReq())
			case *sessionv1.SessionMessage_Game:
				session.GameState.SocketRouter(newSessionEvent.Event.GetGame(), newSessionEvent.Player.ID)
			}

		case <-session.GameState.GetContext().Done():
			session.IsGameEnded = true
			if session.Chat.IsOn {
				expDur := 30 * time.Second
				if session.Bot != nil {
					text := fmt.Sprintf("چت تا %d ثانیه دیگه بسته میشه", int(expDur.Seconds()))
					for _, player := range session.Players {
						_, err := session.Bot.Send(&telebot.User{ID: int64(player.TgID)}, text)
						if err != nil {
							slog.Error("can't send end game chat message", "err", err)
						}
					}
				}
				time.Sleep(expDur)
			}
			session.CleanAndDisconnect(allSession)
			return

		case <-session.CommandNotifire:
			if len(session.Commands) > 0 {
				com := session.PopCommand()
				session.ApplyCommand(com)
			}
		}
	}
}

func (session *GameSession) CleanAndDisconnect(allSession *AllSession) {
	if session.Chat.IsOn {
		text := "چت قطع شد"
		for _, player := range session.Players {
			_, err := session.Bot.Send(&telebot.User{ID: int64(player.TgID)}, text, keybul.WelcomeReplyKeyboard)
			if err != nil {
				slog.Error("can't send chat ended message", "err", err)
			}
		}
	}

	allSession.Mutex.Lock()
	defer allSession.Mutex.Unlock()

	for _, player := range session.Players {
		delete(allSession.Sessions, strconv.Itoa(player.TgID))
	}
}
func (session *GameSession) HandleCallback(c telebot.Context, queries *database.Queries) error {
	callbackData := c.Callback().Data
	if callbackData == "join" {
		personRow, err := queries.GetTgUserByTgID(context.Background(), int(c.Sender().ID))
		if err != nil {
			slog.Error("can not get user at session handle callback", "err", err)
			return c.RespondText("خطا")
		}
		if c.Sender().ID == int64(session.Players[0].TgID) {
			text := "خودت بازیو ساختی تو بازی هستی"
			return c.RespondText(text)
		}

		newJoinCommand := NewJoinSessionCommand(session, NewSessionPlayer(personRow.ID, personRow.TgID, personRow.Name))
		session.PushCommand(newJoinCommand)

		text := "اضافه شدی بازی شروع شد"
		err = c.RespondText(text)
		if err != nil {
			return err
		}
	}
	err := session.GameState.CallBackRouter(c)
	if err != nil {
		slog.Error("error in call back router", "error", err)
		return c.RespondText("خطا")
	}

	return c.RespondText("none")
}
