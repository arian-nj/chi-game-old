package bot

import (
	"fmt"
	"log/slog"
	"strconv"
	"sync"
	"time"

	gamesessions "github.com/arian-nj/chibazi/game_sessions"
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

func (app *BotApplication) AddTicket(gameType gametype.GameType, newTicket *Ticket) {

	queue := app.MatchMaking.WaitingPlayers[gameType]
	app.MatchMaking.WaitingPlayers[gameType] = append(queue, newTicket)

}

func (app *BotApplication) CheckIsAllowedToPlay(playerId int) bool {
	_, isFound := app.AllSessions.Sessions[strconv.Itoa(playerId)]
	return !isFound
}

func (app *BotApplication) MakeMatches() {
	defer app.MatchMaking.Mutex.Unlock()
	var doFlag = false
	for {
		doFlag = false
		for gameTypeKey, ticketsList := range app.MatchMaking.WaitingPlayers {
			app.MatchMaking.Mutex.Lock()
			if len(ticketsList) >= 2 {
				doFlag = true
				ticketOne := ticketsList[0]
				ticketTwo := ticketsList[1]
				app.MatchMaking.WaitingPlayers[gameTypeKey] = ticketsList[2:]
				app.createRandomGame(gameTypeKey, []*Ticket{ticketOne, ticketTwo})
			}
			app.MatchMaking.Mutex.Unlock()
		}
		if !doFlag {
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (app *BotApplication) createRandomGame(gameType gametype.GameType, tickets []*Ticket) {
	playerOne := consoleplayer.NewConsolePlayer(tickets[0].Name, tickets[0].UserID).SetMessageSig(tickets[0].MessageID)
	playerTwo := consoleplayer.NewConsolePlayer(tickets[1].Name, tickets[1].UserID).SetMessageSig(tickets[1].MessageID)

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

	var newSession *gamesessions.GameSession
	switch gameType {
	case gametype.XOGameType3X3, gametype.XOGameType5X5:
		newXOGame := xoconsole.NewXOGame(gametype.XOGameType3X3, app.Queries)
		newXOGame.AddPlayer(playerOne)
		newXOGame.AddPlayer(playerTwo)

		newSession = gamesessions.NewGameSession(app.AllSessions, app.Bot, gameType, newXOGame)
		app.AllSessions.Add(strconv.Itoa(playerOne.TgID), newSession)
		app.AllSessions.Add(strconv.Itoa(playerTwo.TgID), newSession)

		err := newXOGame.StartGame(app.Bot)
		if err != nil {
			slog.Error("error in starting random xo match", "error", err)
		}

	default:
		slog.Error("not possible")
		return
	}

	for _, player := range newSession.GameState.Players() {
		for _, oppPlayer := range newSession.GameState.Players() {
			if player.TgID == oppPlayer.TgID {
				continue
			}
			_, err := app.Bot.Send(&telebot.User{ID: int64(player.TgID)}, FoundOpponentText(oppPlayer.Name), telebot.ModeMarkdownV2, StopChatReplyKeyboard)
			if err != nil {
				slog.Error("can't send found opponent message ", "error", err)
			}
		}
	}
}

func FoundOpponentText(oppName string) string {
	text := ""
	text += "🕹 بازی شروع شد☝️\n\n"
	text += fmt.Sprintf("👀 به حریفت *%s* سلام کن 🤝", oppName)
	return text
}

var StopChatReplyKeyboard = &telebot.ReplyMarkup{
	ReplyKeyboard: [][]telebot.ReplyButton{
		{
			{Text: StopChatButtonText},
		},
	},
	ResizeKeyboard: true,
}

func (app *BotApplication) inlineResultFeedbackHandler(c telebot.Context) error {
	resultId := c.InlineResult().ResultID
	messageID := c.InlineResult().MessageID

	switch gametype.GameType(resultId) {
	case gametype.XOGameType3X3:
		newXOGame := xoconsole.NewXOGame(gametype.XOGameType3X3, app.Queries)
		newGameSession := gamesessions.NewGameSession(app.AllSessions, app.Bot, gametype.XOGameType3X3, newXOGame)

		app.AllSessions.Add(messageID, newGameSession)

		newXOGame.ViaMessageId = messageID
		return newXOGame.SendJoinPanel(c)

	case gametype.XOGameType5X5:
		newXOGame := xoconsole.NewXOGame(gametype.XOGameType5X5, app.Queries)
		newGameSession := gamesessions.NewGameSession(app.AllSessions, app.Bot, gametype.XOGameType3X3, newXOGame)

		app.AllSessions.Add(messageID, newGameSession)

		newXOGame.ViaMessageId = messageID
		return newXOGame.SendJoinPanel(c)
	}
	return c.RespondAlert("این بازیرو ندارم!")
}

func (app *BotApplication) RemovePlayerFromMatchMaking(userID int) bool {

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
