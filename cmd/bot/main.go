package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/arian-nj/chibazi/api"
	"github.com/arian-nj/chibazi/bot"
	"github.com/arian-nj/chibazi/database"
	"github.com/arian-nj/chibazi/db"
	gamesessions "github.com/arian-nj/chibazi/game_sessions"
	xoconsole "github.com/arian-nj/chibazi/games/xo_console"
	"github.com/arian-nj/chibazi/internals/config"
	gametype "github.com/arian-nj/chibazi/internals/game_type"
	matchmaking "github.com/arian-nj/chibazi/match_making"
	"github.com/jackc/pgx/v5/pgxpool"
	"gopkg.in/telebot.v4"
)

// NOTE: use matchmaking as a package not to run code u fucking idiot
type GlobalVars struct {
	Config      *config.Config
	Queries     *database.Queries
	Conn        *pgxpool.Pool
	Wg          *sync.WaitGroup
	AllSessions *gamesessions.AllSession
	MatchMaking *matchmaking.MatchMaking
	Bot         *telebot.Bot
}

func NewGlobalVars() *GlobalVars {
	return &GlobalVars{
		AllSessions: &gamesessions.AllSession{
			Sessions: map[string]*gamesessions.GameSession{},
			Mutex:    sync.Mutex{},
		},
		Wg: &sync.WaitGroup{},
	}
}

func main() {
	GlobalVars := NewGlobalVars()

	var err error
	GlobalVars.Config, err = config.ParseConfig()
	if err != nil {
		panic(err)
	}

	err = db.Migrate(GlobalVars.Config.DatabseUrl)
	if err != nil {
		slog.Error("Failed to migrate database", "err", err)
		return
	}

	GlobalVars.Conn, err = pgxpool.New(context.Background(), GlobalVars.Config.DatabseUrl)
	if err != nil {
		slog.Error("can not make a new connection ", "err", err)
		return
	}
	defer GlobalVars.Conn.Close()

	GlobalVars.Queries = database.New(GlobalVars.Conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	err = GlobalVars.Conn.Ping(ctx)
	if err != nil {
		slog.Error("Failed to connect to Database", "err", err)
		cancel()
		return
	}
	GlobalVars.MatchMaking = matchmaking.NewMatchMaking(GlobalVars.AllSessions, GlobalVars.Queries)

	slog.Info("Connected to Database")

	parentCtx, cancel := context.WithCancel(context.Background())

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	app := api.NewApiApplication(GlobalVars.Config, GlobalVars.Queries, GlobalVars.AllSessions, GlobalVars.MatchMaking)
	go app.RunApi(parentCtx, GlobalVars.Wg)

	botApp := bot.NewBotApplication(GlobalVars.Config, GlobalVars.Queries, GlobalVars.AllSessions, GlobalVars.MatchMaking)
	bot, err := botApp.MakeBot()
	if err != nil {
		slog.Error("Failed to make bot", "err", err)
		return
	}
	GlobalVars.Bot = bot
	go botApp.RunBot(bot, parentCtx, GlobalVars.Wg)

	go ClearDeadGamesCron(GlobalVars.AllSessions)

	go GlobalVars.MakeMatches()

	<-quit
	cancel()
	GlobalVars.Wg.Wait()
}

func ClearDeadGamesCron(allSessions *gamesessions.AllSession) {
	for {
		nowTime := time.Now()
		for key, gameSession := range allSessions.Sessions {
			if nowTime.Sub(gameSession.CreatedAt) > gameSession.ExpireDuaration {
				allSessions.Mutex.Lock()
				delete(allSessions.Sessions, key)
				allSessions.Mutex.Unlock()
			}
		}
		time.Sleep(1 * time.Minute)
	}
}

func (gv *GlobalVars) MakeMatches() {
	defer gv.MatchMaking.Mutex.Unlock()
	var doFlag = false
	for {
		doFlag = false
		for gameTypeKey, ticketsList := range gv.MatchMaking.WaitingPlayers {
			gv.MatchMaking.Mutex.Lock()
			if len(ticketsList) >= 2 {
				doFlag = true
				ticketOne := ticketsList[0]
				ticketTwo := ticketsList[1]
				gv.MatchMaking.WaitingPlayers[gameTypeKey] = ticketsList[2:]
				gv.createRandomGame(gameTypeKey, []*matchmaking.Ticket{ticketOne, ticketTwo})
			}
			gv.MatchMaking.Mutex.Unlock()
		}
		if !doFlag {
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (gv *GlobalVars) createRandomGame(gameType gametype.GameType, tickets []*matchmaking.Ticket) {
	var newSession *gamesessions.GameSession

	switch gameType {

	case gametype.XOGameType3X3, gametype.XOGameType5X5:
		newXOGame := xoconsole.NewXOGame(gametype.XOGameType3X3, gv.Queries)
		newSession = gamesessions.NewGameSession(gv.AllSessions, gv.Bot, gameType, newXOGame)

	default:
		slog.Error("not possible random game")
		return

	}

	playerOne := gamesessions.NewSessionPlayer(tickets[0].TgID, tickets[0].Name)
	playerTwo := gamesessions.NewSessionPlayer(tickets[1].TgID, tickets[1].Name)

	newSession.AddPlayer(playerOne)
	newSession.AddPlayer(playerTwo)

	gv.AllSessions.Add(strconv.Itoa(playerOne.TgID), newSession)
	gv.AllSessions.Add(strconv.Itoa(playerTwo.TgID), newSession)

	err := newSession.StartGame()
	if err != nil {
		slog.Error("error in starting random xo match", "error", err)
		return
	}
	_, err = gv.Queries.CreateGameSession(context.Background())

	for _, ticket := range tickets {
		ticket.MatchFound <- newSession
	}

	bot.SendFoundOpponentMessage(newSession.GameState.Players(), gv.Bot)
}
