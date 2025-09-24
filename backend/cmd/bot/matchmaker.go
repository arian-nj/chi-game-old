package main

import (
	"context"
	"log/slog"
	_ "net/http/pprof"
	"strconv"
	"time"

	gamesessions "github.com/arian-nj/chibazi/backend/game_sessions"
	"github.com/arian-nj/chibazi/backend/games/conn4"
	"github.com/arian-nj/chibazi/backend/games/games"
	"github.com/arian-nj/chibazi/backend/games/xo"
	"github.com/arian-nj/chibazi/backend/internals/utils"
	matchmaking "github.com/arian-nj/chibazi/backend/match_making"
)

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
				gv.createRandomGame(gameTypeKey, ticketOne, ticketTwo)
			}
			gv.MatchMaking.Mutex.Unlock()
		}
		if !doFlag {
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (gv *GlobalVars) createRandomGame(gameType games.GameType, ticketOne *matchmaking.Ticket, ticketTwo *matchmaking.Ticket) {
	newSessionRow, err := gv.Queries.CreateSession(context.Background(), string(gamesessions.RandomSession))
	if err != nil {
		slog.Error("can't create new random game", "error", err)
	}
	newGameSession := gamesessions.NewGameSession(gv.Bot, gv.Queries, newSessionRow.ID, gv.AllSessions)
	newGameSession.Subscribe(gamesessions.NewSessionTelegramBotListener(ticketOne.UserID, ticketOne.TgID, gv.Bot, ""))
	newGameSession.Subscribe(gamesessions.NewSessionTelegramBotListener(ticketTwo.UserID, ticketTwo.TgID, gv.Bot, ""))

	var newGame games.Game

	switch gameType {
	case games.XOGameType3X3, games.XOGameType5X5:
		newXoGame := xo.NewXOGame(newGameSession.SessionCtx, games.XOGameType3X3, gv.Queries)
		newGame = newXoGame
	case games.Conn4GameType:
		newConn4Game := conn4.NewConn4State(newGameSession.SessionCtx, games.Conn4GameType, gv.Queries)
		newGame = newConn4Game

	default:
		slog.Error("not possible random game")
		return
	}
	newGameSession.GameState = newGame
	newGameSession.RunBgMonitor()

	playerOne := gamesessions.NewSessionPlayer(ticketOne.UserID, ticketOne.TgID, ticketOne.Name)
	playerTwo := gamesessions.NewSessionPlayer(ticketTwo.UserID, ticketTwo.TgID, ticketTwo.Name)

	newGameSession.AddPlayerToSession(playerOne)
	newGameSession.AddPlayerToSession(playerTwo)

	gv.AllSessions.Add(strconv.Itoa(playerOne.ID), newGameSession)
	gv.AllSessions.Add(strconv.Itoa(playerTwo.ID), newGameSession)

	for _, ticket := range []*matchmaking.Ticket{ticketOne, ticketTwo} {
		ticket.MatchFoundChan <- newGameSession
	}

	utils.RunBackgroundTask(func() {
		err = newGameSession.StartGame()
		if err != nil {
			slog.Error("error in starting random xo match", "error", err)
			return
		}
		gamesessions.SendFoundOpponentMessage(newGameSession.Players, gv.Bot)

	})
}
