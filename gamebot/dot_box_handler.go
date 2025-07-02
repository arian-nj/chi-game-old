package gamebot

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	dotbox "github.com/arian-nj/ultrun/internals/dot_box"
	"github.com/arian-nj/ultrun/internals/random"
	"gopkg.in/telebot.v4"
)

type DotBoxGame struct { // of GameInterface type
	DotBoxBoard        *dotbox.DotBoxBoard
	Hub                *Hub
	CurrentPlayerIndex int

	PlayerOnScore  int
	PlayerTwoScore int
}

func NewDotBoxGame() *DotBoxGame {
	return &DotBoxGame{
		DotBoxBoard: dotbox.NewDotBoard(),
	}
}

const (
	dotBoxStartText = `❌ *نقطه بازی* ⭕️`
)

func (game *DotBoxGame) NextPlayer() {
	if game.CurrentPlayerIndex == len(game.Hub.Players)-1 {
		game.CurrentPlayerIndex = 0
		return
	}
	game.CurrentPlayerIndex += 1
}

func (app *Application) dotBoxInlineResultReciever(c telebot.Context) error {
	sender := c.InlineResult().Sender

	dotBoxGame := NewDotBoxGame()
	hub := NewHub(dotBoxGame, c.InlineResult().MessageID)
	dotBoxGame.Hub = hub
	app.Lobby.Hubs.AddHub(hub)
	hub.AddPlayer(NewHumanPlayer(int(sender.ID), sender.FirstName))

	textMessage := dotBoxStartText + "\n\n🕹 بازیکن " + fmt.Sprintf("[%s](tg://user?id=%d)", sender.FirstName, sender.ID) + " منتظر حریفه"

	_, err := c.Bot().Edit(c.InlineResult(), textMessage, app.JoinGameKeyboard(), telebot.ModeMarkdownV2)

	return err
}

func (game *DotBoxGame) StartGame(app *Application) {
	randIndex := random.GenerateRandomNumber(2)
	game.CurrentPlayerIndex = randIndex

	_, err := app.Bot.Edit(game.Hub, ticStartText, telebot.ModeMarkdownV2,
		CreateInlineKeyboard(
			CreateBotNameInlineButton(app.Bot),
			CreateDotBoxInlineButton(game.DotBoxBoard),
			game.CreatePlayersInlineButton(game.Hub.Players, game.CurrentPlayerIndex),
		),
	)
	if err != nil {
		slog.Error("error when starting dot box game" + err.Error())
		return
	}
}

func (game *DotBoxGame) GetMaxPlayer() int {
	return 2
}
func (game *DotBoxGame) GetGameType() GameType {
	return TicTacToeGameType
}

func getSymbol(r, c int) string {

	if r == 0 && c == dotbox.MaxCellSize-1 {
		return "┐"
	} else if r == 0 && c == 0 {
		return "┌"
	} else if r == dotbox.MaxCellSize-1 && c == 0 {
		return "└"
	} else if r == dotbox.MaxCellSize-1 && c == dotbox.MaxCellSize-1 {
		return "┘"
	}

	if r == 0 {
		return "┬"
	} else if r == dotbox.MaxCellSize-1 {
		return "┴"
	} else if c == 0 {
		return "├"
	} else if c == dotbox.MaxCellSize-1 {
		return "┤"
	}

	return "┿"
}

func CreateDotBoxInlineButton(board *dotbox.DotBoxBoard) [][]telebot.InlineButton {
	buttons := [][]telebot.InlineButton{}

	for r := range dotbox.MaxCellSize {
		row := []telebot.InlineButton{}
		for c := range dotbox.MaxCellSize {
			value := " "
			if board.Board[r][c] != dotbox.Empty {
				value = getSymbol(r, c)
			}
			btn := telebot.InlineButton{
				Text: value,
				Data: string(DotBoxGameType) + "_play_" + strconv.Itoa(r) + strconv.Itoa(c),
			}
			row = append(row, btn)
		}
		buttons = append(buttons, row)
	}
	return buttons
}

func (app *Application) DotBoxCallbackHandlers(c telebot.Context) error {
	callback := c.Callback()
	callbackData := strings.TrimPrefix(callback.Data, string(DotBoxGameType)+"_")
	messageId := c.Callback().MessageID

	hub, is_found := app.Lobby.Hubs[messageId]
	if !is_found {
		return c.RespondAlert("این بازی وجود نداره!")
	}

	dotBoxGame, ok := hub.Game.(*DotBoxGame)
	if !ok {
		panic("dot box callback handler can't convert game interface to dot box game struct")
	}
	if strings.HasPrefix(callbackData, "play_") {
		return dotBoxGame.PlayHandler(c, app)
	}
	return nil
}

func (game *DotBoxGame) PlayHandler(c telebot.Context, app *Application) error {
	callback := c.Callback()
	callbackData := strings.TrimPrefix(callback.Data, string(DotBoxGameType)+"_play_")

	if len(callbackData) != 2 {
		return fmt.Errorf("invalid ttt_play data")
	}

	if game.Hub.Players[game.CurrentPlayerIndex].TgID != int(callback.Sender.ID) {
		return c.RespondText("نوبت تو نیست!")
	}

	xySlice := strings.Split(callbackData, "")
	rstr, cstr := xySlice[0], xySlice[1]
	rint, xerr := strconv.Atoi(rstr)
	cint, yerr := strconv.Atoi(cstr)
	if xerr != nil || yerr != nil {
		c.RespondAlert("یه مشکلی هست")
	}

	moveType := dotbox.Empty
	if game.CurrentPlayerIndex == 0 {
		moveType = dotbox.Blue
	} else {
		moveType = dotbox.Red
	}
	isValid, errMsg, isScore := game.DotBoxBoard.PlayMove(rint, cint, moveType)
	if !isValid {
		return c.RespondText(errMsg)
	}

	if isScore {
		if game.CurrentPlayerIndex == 0 {
			game.PlayerOnScore += 1
		} else {
			game.PlayerTwoScore += 2
		}

		err := c.RespondText("امتیاز گرفتی بازم نوبتته")
		if err != nil {
			return err
		}
	} else {
		game.NextPlayer()
	}

	_, err := app.Bot.Edit(game.Hub, dotBoxStartText, telebot.ModeMarkdownV2,
		CreateInlineKeyboard(
			CreateBotNameInlineButton(app.Bot),
			CreateDotBoxInlineButton(game.DotBoxBoard),
			game.CreatePlayersInlineButton(game.Hub.Players, game.CurrentPlayerIndex),
		),
	)

	if err != nil {
		return err
	}
	return nil
}

func (game *DotBoxGame) CreatePlayersInlineButton(humanPlayers []*HumanPlayer, CurrentPlayerTurn int) [][]telebot.InlineButton {
	buttons := make([][]telebot.InlineButton, 0)
	for index, hplayer := range humanPlayers {

		yourTurn := ""
		if CurrentPlayerTurn == index {
			yourTurn = "🎮"
		}

		// emoji := ""
		score := fmt.Sprint(game.PlayerTwoScore)
		if index == 0 {
			score = fmt.Sprint(game.PlayerOnScore)
		}

		name := hplayer.Name
		if len(name) > 20 {
			name = name[:20] + "..."
		}

		row := make([]telebot.InlineButton, 2)
		row = append(row, telebot.InlineButton{
			Text: fmt.Sprintf("%s %s (%s)", yourTurn, name, score),
			URL:  fmt.Sprintf("tg://user?id=%d", hplayer.TgID),
		})
		buttons = append(buttons, row)
	}

	return buttons
}
