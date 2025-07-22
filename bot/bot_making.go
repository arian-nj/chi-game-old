package bot

import (
	"fmt"
	gamesessions "github.com/arian-nj/chibazi/game_sessions"
	xoconsole "github.com/arian-nj/chibazi/games/xo_console"
	gametype "github.com/arian-nj/chibazi/internals/game_type"
	"gopkg.in/telebot.v4"
)

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
