package gamebot

import (
	"log/slog"
	"strconv"

	xoconsole "github.com/arian-nj/chibazi/games/xo_console"
	gametype "github.com/arian-nj/chibazi/internals/game_type"
	matchmaking "github.com/arian-nj/chibazi/internals/match_making"
	"gopkg.in/telebot.v4"
)

func (app *Application) MakeMatches() {
	for {
		for gameTypeKey, ticketsList := range app.MatchMaking.WaitingPlayers {
			app.MatchMaking.Mutex.Lock()
			if len(ticketsList) >= 2 {
				ticketOne := ticketsList[0]
				ticketTwo := ticketsList[1]
				app.MatchMaking.WaitingPlayers[gameTypeKey] = ticketsList[2:]
				app.createRandomGame(gameTypeKey, []*matchmaking.Ticket{ticketOne, ticketTwo})
			}
			app.MatchMaking.Mutex.Unlock()
		}
	}
}

func (app *Application) createRandomGame(gameType gametype.GameType, tickets []*matchmaking.Ticket) {
	switch gameType {
	case gametype.XOGameType3X3:
		playerOne := xoconsole.NewPlayer(tickets[0].Name, tickets[0].UserID).SetMessageSig(tickets[0].MessageID)
		playerTwo := xoconsole.NewPlayer(tickets[1].Name, tickets[1].UserID).SetMessageSig(tickets[1].MessageID)
		newXOGame := xoconsole.NewXOGame(gametype.XOGameType3X3, app.Queries)
		newXOGame.AddPlayer(playerOne)
		newXOGame.AddPlayer(playerTwo)

		app.Lobby.XOGames[strconv.Itoa(playerOne.TgID)] = newXOGame
		app.Lobby.XOGames[strconv.Itoa(playerTwo.TgID)] = newXOGame

		err := newXOGame.StartGame(app.Bot)
		if err != nil {
			slog.Error("error in starting random xo match", "error", err)
		}

	}
}

func (app *Application) inlineResultFeedbackHandler(c telebot.Context) error {
	resultId := c.InlineResult().ResultID
	messageID := c.InlineResult().MessageID

	switch gametype.GameType(resultId) {
	case gametype.XOGameType3X3:
		newXOGame := xoconsole.NewXOGame(gametype.XOGameType3X3, app.Queries)
		app.Lobby.XOGames[messageID] = newXOGame
		newXOGame.ViaMessageId = messageID
		return newXOGame.SendJoinPanel(c)

	case gametype.XOGameType5X5:
		newXOGame := xoconsole.NewXOGame(gametype.XOGameType5X5, app.Queries)
		app.Lobby.XOGames[messageID] = newXOGame
		newXOGame.ViaMessageId = messageID
		return newXOGame.SendJoinPanel(c)
	}
	return c.RespondAlert("این بازیرو ندارم!")
}
