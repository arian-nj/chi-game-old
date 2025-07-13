package bot

import (
	"log/slog"
	"strconv"
	"sync"
	"time"

	xoconsole "github.com/arian-nj/chibazi/games/xo_console"
	consoleplayer "github.com/arian-nj/chibazi/internals/console_player"
	gametype "github.com/arian-nj/chibazi/internals/game_type"
	"gopkg.in/telebot.v4"
)

type MatchMaking struct {
	WaitingPlayers map[gametype.GameType][]*Ticket
	Mutex          sync.Mutex
}

type Ticket struct {
	Name      string
	UserID    int
	MessageID int
	GameType  gametype.GameType
	Timestamp time.Time
}

func NewTicket(name string, userID, messageID int, gameType gametype.GameType) *Ticket {
	return &Ticket{
		UserID:    userID,
		Name:      name,
		MessageID: messageID,
		GameType:  gameType,
	}
}

func (app *Application) MakeMatches() {
	for {
		for gameTypeKey, ticketsList := range app.MatchMaking.WaitingPlayers {
			app.MatchMaking.Mutex.Lock()
			if len(ticketsList) >= 2 {
				ticketOne := ticketsList[0]
				ticketTwo := ticketsList[1]
				app.MatchMaking.WaitingPlayers[gameTypeKey] = ticketsList[2:]
				app.createRandomGame(gameTypeKey, []*Ticket{ticketOne, ticketTwo})
			}
			app.MatchMaking.Mutex.Unlock()
		}
	}
}

func (app *Application) createRandomGame(gameType gametype.GameType, tickets []*Ticket) {
	playerOne := consoleplayer.NewPlayer(tickets[0].Name, tickets[0].UserID).SetMessageSig(tickets[0].MessageID)
	playerTwo := consoleplayer.NewPlayer(tickets[1].Name, tickets[1].UserID).SetMessageSig(tickets[1].MessageID)

	for _, player := range []*consoleplayer.ConsolePlayer{playerOne, playerTwo} {

		go func(inPlayer consoleplayer.ConsolePlayer) {
			err := app.Bot.Delete(&inPlayer)
			if err != nil {
				slog.Error("can't delete match making message", "err", err)
				return
			}
		}(*player)

		msg, err := app.Bot.Send(&telebot.User{ID: int64(player.TgID)}, "بازی")
		if err != nil {
			slog.Error("can't send base game message", "err", err)
		}
		player.MessageID = msg.ID
	}

	var newSession GameSession
	switch gameType {
	case gametype.XOGameType3X3:
		newXOGame := xoconsole.NewXOGame(gametype.XOGameType3X3, app.Queries)
		newXOGame.AddPlayer(playerOne)
		newXOGame.AddPlayer(playerTwo)

		newSession = NewGameSession(gameType, newXOGame)
		app.AddGameSession(strconv.Itoa(playerOne.TgID), newSession)
		app.AddGameSession(strconv.Itoa(playerTwo.TgID), newSession)

		err := newXOGame.StartGame(app.Bot)
		if err != nil {
			slog.Error("error in starting random xo match", "error", err)
		}
	default:
		slog.Error("not possible")
		return
	}
	// for _,player := range newSession.GameState.Players() {

	// }
}

func (app *Application) inlineResultFeedbackHandler(c telebot.Context) error {
	resultId := c.InlineResult().ResultID
	messageID := c.InlineResult().MessageID

	switch gametype.GameType(resultId) {
	case gametype.XOGameType3X3:
		newXOGame := xoconsole.NewXOGame(gametype.XOGameType3X3, app.Queries)
		newGameSession := NewGameSession(gametype.XOGameType3X3, newXOGame)

		app.AddGameSession(messageID, newGameSession)

		newXOGame.ViaMessageId = messageID
		return newXOGame.SendJoinPanel(c)

	case gametype.XOGameType5X5:
		newXOGame := xoconsole.NewXOGame(gametype.XOGameType5X5, app.Queries)
		newGameSession := NewGameSession(gametype.XOGameType3X3, newXOGame)

		app.AddGameSession(messageID, newGameSession)

		newXOGame.ViaMessageId = messageID
		return newXOGame.SendJoinPanel(c)
	}
	return c.RespondAlert("این بازیرو ندارم!")
}

func (app *Application) RemovePlayerFromMatchMaking(userID int) bool {
	app.MatchMaking.Mutex.Lock()
	defer app.MatchMaking.Mutex.Unlock()

	for gameType, tickets := range app.MatchMaking.WaitingPlayers {
		for index, ticket := range tickets {
			if ticket.UserID == int(userID) {
				li := app.MatchMaking.WaitingPlayers[gameType]
				app.MatchMaking.WaitingPlayers[gameType] = append(li[:index], li[index+1:]...)
				return true
			}
		}
	}
	return false
}
