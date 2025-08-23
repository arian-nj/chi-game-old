package matchmaking

import (
	"sync"
	"time"

	"github.com/arian-nj/chibazi/database"
	gamesessions "github.com/arian-nj/chibazi/game_sessions"
	gametype "github.com/arian-nj/chibazi/internals/game_type"
)

type MatchMaking struct {
	WaitingPlayers map[gametype.GameType][]*Ticket
	Mutex          sync.Mutex
	AllSessions    *gamesessions.AllSession
	Queries        *database.Queries
}

func NewMatchMaking(allSessions *gamesessions.AllSession, queries *database.Queries) *MatchMaking {
	return &MatchMaking{
		AllSessions:    allSessions,
		Queries:        queries,
		WaitingPlayers: map[gametype.GameType][]*Ticket{},
	}
}

type PlatformType string

const (
	TelegramConsole PlatformType = "telegram_console"
	TelegramMiniApp PlatformType = "telegram_miniapp"
)

type Ticket struct {
	Name   string
	UserID int
	TgID   int

	Platform PlatformType
	GameType gametype.GameType

	MatchFoundChan chan *gamesessions.GameSession

	Timestamp time.Time
}

func NewTicket(name string, userID int, tgID int, gameType gametype.GameType) *Ticket {
	return &Ticket{
		UserID:         userID,
		TgID:           tgID,
		Name:           name,
		GameType:       gameType,
		MatchFoundChan: make(chan *gamesessions.GameSession, 1),
	}
}

func (mm *MatchMaking) HasTicket(tgUserID int) bool {
	mm.Mutex.Lock()
	defer mm.Mutex.Unlock()

	for _, tickets := range mm.WaitingPlayers {
		for _, ticket := range tickets {
			if ticket.TgID == tgUserID {
				return true
			}
		}
	}
	return false

}

func (mm *MatchMaking) PushTicket(newTicket *Ticket) {
	mm.Mutex.Lock()
	defer mm.Mutex.Unlock()

	queue := mm.WaitingPlayers[newTicket.GameType]
	mm.WaitingPlayers[newTicket.GameType] = append(queue, newTicket)
}
func (mm *MatchMaking) RemovePlayerTicket(tgUserID int) bool {
	mm.Mutex.Lock()
	defer mm.Mutex.Unlock()

	for gameType, tickets := range mm.WaitingPlayers {
		for index, ticket := range tickets {
			if ticket.TgID == tgUserID {
				li := mm.WaitingPlayers[gameType]
				mm.WaitingPlayers[gameType] = append(li[:index], li[index+1:]...)
				return true
			}
		}
	}
	return false
}

//
// func (mm *MatchMaking) MakeMatches() {
// 	defer mm.Mutex.Unlock()
// 	var doFlag = false
// 	for {
// 		doFlag = false
// 		for gameTypeKey, ticketsList := range mm.WaitingPlayers {
// 			mm.Mutex.Lock()
// 			if len(ticketsList) >= 2 {
// 				doFlag = true
// 				ticketOne := ticketsList[0]
// 				ticketTwo := ticketsList[1]
// 				mm.WaitingPlayers[gameTypeKey] = ticketsList[2:]
// 				mm.createRandomGame(gameTypeKey, []*Ticket{ticketOne, ticketTwo})
// 			}
// 			mm.Mutex.Unlock()
// 		}
// 		if !doFlag {
// 			time.Sleep(100 * time.Millisecond)
// 		}
// 	}
// }
// func (mm *MatchMaking) createRandomGame(gameType gametype.GameType, tickets []*Ticket) {
// 	playerOne := consoleplayer.NewConsolePlayer(tickets[0].Name, tickets[0].UserID).SetMessageSig(tickets[0].MessageID)
// 	playerTwo := consoleplayer.NewConsolePlayer(tickets[1].Name, tickets[1].UserID).SetMessageSig(tickets[1].MessageID)
// 	ClearMatchmakingMessage([]*consoleplayer.ConsolePlayer{playerOne, playerTwo}, mm.Bot)
//
// 	var newSession *gamesessions.GameSession
// 	switch gameType {
// 	case gametype.XOGameType3X3, gametype.XOGameType5X5:
// 		newXOGame := xoconsole.NewXOGame(gametype.XOGameType3X3, mm.Queries)
// 		newXOGame.AddPlayer(playerOne)
// 		newXOGame.AddPlayer(playerTwo)
//
// 		newSession = gamesessions.NewGameSession(mm.AllSessions, mm.Bot, gameType, newXOGame)
// 		mm.AllSessions.Add(strconv.Itoa(playerOne.TgID), newSession)
// 		mm.AllSessions.Add(strconv.Itoa(playerTwo.TgID), newSession)
//
// 		err := newXOGame.StartGame(mm.Bot)
// 		if err != nil {
// 			slog.Error("error in starting random xo match", "error", err)
// 		}
//
// 	default:
// 		slog.Error("not possible")
// 		return
// 	}
// 	SendFoundOpponentMessage(newSession.GameState.Players(), mm.Bot)
//
// }
//
// func ClearMatchmakingMessage(players []*consoleplayer.ConsolePlayer, bot *telebot.Bot) {
// 	for _, player := range players {
//
// 		go func(inPlayer consoleplayer.ConsolePlayer) {
// 			err := bot.Delete(&inPlayer)
// 			if err != nil {
// 				slog.Error("can't delete match making message", "err", err)
// 				return
// 			}
// 		}(*player)
//
// 		msg, err := bot.Send(&telebot.User{ID: int64(player.TgID)}, "بازی")
// 		if err != nil {
// 			slog.Error("can't send base game message", "err", err)
// 		}
// 		player.MessageID = msg.ID
// 	}
// }
//
// // SendFoundOpponentMessage sends a message to all players that they have found an opponent
// func SendFoundOpponentMessage(players []*consoleplayer.ConsolePlayer, bot *telebot.Bot) {
// 	for _, player := range players {
// 		for _, oppPlayer := range players {
// 			if player.TgID == oppPlayer.TgID {
// 				continue
// 			}
// 			_, err := bot.Send(&telebot.User{ID: int64(player.TgID)}, FoundOpponentText(oppPlayer.Name), telebot.ModeMarkdownV2, StopChatReplyKeyboard)
// 			if err != nil {
// 				slog.Error("can't send found opponent message ", "error", err)
// 			}
// 		}
// 	}
// }
