package gamesessions

import (
	"context"
	"errors"
	"log/slog"
	"strconv"
	"time"

	"github.com/arian-nj/chibazi/backend/database"
	"github.com/arian-nj/chibazi/backend/games/games"
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

const ExpirationDur = 30 * time.Second

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
	GameState   games.Game

	MsgChnl chan *SessionEvent

	Players []*SessionPlayer

	CancelSession context.CancelFunc
	SessionCtx    context.Context

	CreatedAt       time.Time
	ExpireDuaration time.Duration

	allSession *AllSession

	*commander.Commander

	ShutdownTimer <-chan time.Time
}

func NewGameSession(bot *telebot.Bot, Queries *database.Queries, sessionId int, allSession *AllSession) *GameSession {
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

		Commander:     commander.NewCommander(),
		allSession:    allSession,
		ShutdownTimer: make(<-chan time.Time),
	}
	return gs
}

func (session *GameSession) RunBgMonitor() {
	utils.RunBackgroundTask(func() {
		session.MonitorGameSession()
	})
}

func (session *GameSession) StartGame() error {
	if session.GameState == nil {
		return errors.New("failed to start game field GameState is nil")
	}
	// session.PushCommand(
	// 	NewGameStartSessionCommand(session),
	// )

	for _, player := range session.Players {
		session.GameState.AddPlayer(player.ID, player.Name, player.TgID, player.Socket)
	}

	for _, suber := range session.Subscribers {
		if listener, ok := suber.(*SessionTelegramBotListener); ok {
			session.GameState.SubToTelegram(listener.UserID, listener.Bot, "")
		}
		if listener, ok := suber.(*SessionTelegramViaListener); ok {
			session.GameState.SubToTelegram(0, listener.Bot, listener.ViaMessageId)
		}
	}

	session.GameState.OnEnd(session.EndGame)
	return session.GameState.StartGame()
}

func (session *GameSession) AddPlayerToSession(player *SessionPlayer) {
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

func (session *GameSession) EndGame() {
	session.IsGameEnded = true

	gameEnded := NewGameEndedSessionCommand(session)
	session.PushCommand(gameEnded)

	if session.Chat.IsOn == false {
		return
	}
	session.ShutdownTimer = time.After(30 * time.Second)

}

func (session *GameSession) MonitorGameSession() {
	for {
		select {
		case newSessionEvent := <-session.MsgChnl:
			switch newSessionEvent.Event.Content.(type) {
			case *sessionv1.SessionMessage_ChatReq:
				session.SocketRequestSendMsg(newSessionEvent.Player, newSessionEvent.Event.GetChatReq())
			case *sessionv1.SessionMessage_Game:
				session.GameState.SocketRouter(newSessionEvent.Event.GetGame(), newSessionEvent.Player.ID)
			}
		case <-session.CommandNotifire:
			if len(session.Commands) > 0 {
				com := session.PopCommand()
				session.ApplyCommand(com)
			}
		case <-session.ShutdownTimer:
			session.CleanAndDisconnect()
			return
		}
	}
}

func (session *GameSession) CleanAndDisconnect() {
	if session.Chat.IsOn {
		text := "چت قطع شد"
		for _, player := range session.Players {
			_, err := session.Bot.Send(&telebot.User{ID: int64(player.TgID)}, text, keybul.WelcomeReplyKeyboard)
			if err != nil {
				slog.Error("can't send chat ended message", "err", err)
			}
		}
	}

	session.allSession.Mutex.Lock()
	defer session.allSession.Mutex.Unlock()

	for _, player := range session.Players {
		delete(session.allSession.Sessions, strconv.Itoa(player.TgID))
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
		return nil
	}

	err := session.GameState.CallBackRouter(c)
	if err != nil {
		slog.Error("error in call back router", "error", err)
		return c.RespondText("خطا")
	}
	return nil
}
