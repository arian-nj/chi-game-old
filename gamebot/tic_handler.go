package gamebot

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/arian-nj/ultrun/internals/random"
	tictactoe "github.com/arian-nj/ultrun/internals/tic-tac-toe"
	"gopkg.in/telebot.v4"
)

type TicGame struct { // of GameInterface type
	TicBoard           *tictactoe.TicBoard
	Hub                *Hub
	CurrentPlayerIndex int
}

func (game *TicGame) NextPlayer() {
	if game.CurrentPlayerIndex == len(game.Hub.Players)-1 {
		game.CurrentPlayerIndex = 0
		return
	}
	game.CurrentPlayerIndex += 1
}

func NewTicGame() *TicGame {
	return &TicGame{
		TicBoard: tictactoe.NewTicBoard(),
	}
}

func (game *TicGame) StartGame(app *Application) {
	randIndex := random.GenerateRandomNumber(2)
	game.CurrentPlayerIndex = randIndex

	_, err := app.Bot.Edit(game.Hub, ticStartText, telebot.ModeMarkdownV2,
		CreateInlineKeyboard(
			CreateTicBoardInlineButton(game.TicBoard),
			CreatePlayersInlineButton(game.Hub.Players, game.CurrentPlayerIndex),
		),
	)
	if err != nil {
		slog.Error("error when starting game", err.Error())
		return
	}
}

func (game *TicGame) JoinGame(app *Application, player *HumanPlayer) bool {
	if len(game.Hub.Players) >= 2 {
		return false
	}
	game.Hub.AddPlayer(player)
	if len(game.Hub.Players) == 2 {
		game.Hub.Game.StartGame(app)
	}
	return true
}

func (game *TicGame) EndGame() {}

func (game *TicGame) GetGameType() GameType {
	return TicTacToe
}

func (app *Application) ticInlineResultReciever(c telebot.Context) error {
	sender := c.InlineResult().Sender

	ticGame := NewTicGame()
	hub := NewHub("دوز بازی", ticGame, c.InlineResult().MessageID)
	ticGame.Hub = hub
	app.Lobby.Hubs.AddHub(hub)
	hub.AddPlayer(NewHumanPlayer(int(sender.ID), sender.FirstName))

	textMessage := ticStartText + "\n\n🕹 بازیکن " + fmt.Sprintf("[%s](tg://user?id=%d)", sender.FirstName, sender.ID) + " منتظر حریفه"

	_, err := c.Bot().Edit(c.InlineResult(), textMessage, app.JoinGameKeyboard(), telebot.ModeMarkdownV2)

	return err
}

func CreateInlineKeyboard(buttonGroups ...[][]telebot.InlineButton) *telebot.ReplyMarkup {
	var inlineKeyboard [][]telebot.InlineButton

	for _, group := range buttonGroups {
		inlineKeyboard = append(inlineKeyboard, group...)
	}

	return &telebot.ReplyMarkup{
		InlineKeyboard: inlineKeyboard,
	}
}

func (app *Application) XOCallbackHandlers(c telebot.Context) error {
	callback := c.Callback()
	callbackData := strings.TrimPrefix(callback.Data, "xo_")
	messageId := c.Callback().MessageID

	hub, is_found := app.Lobby.Hubs[messageId]
	if !is_found {
		return c.RespondAlert("این بازی وجود نداره!")
	}

	ticGame, ok := hub.Game.(*TicGame)
	if !ok {
		panic("xo callback handler can't convert game interface to ticgame struct")
	}
	if strings.HasPrefix(callbackData, "play_") {
		return ticGame.TicPlayHandler(c, app)
	}
	return nil
}

func (game *TicGame) TicPlayHandler(c telebot.Context, app *Application) error {
	callback := c.Callback()
	callbackData := strings.TrimPrefix(callback.Data, "xo_play_")

	if len(callbackData) != 2 {
		return fmt.Errorf("invalid xo_play data")
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

	moveType := tictactoe.Empty
	if game.CurrentPlayerIndex == 0 {
		moveType = tictactoe.X
	} else {
		moveType = tictactoe.O
	}

	isValid, errMsg := game.TicBoard.PlayMove(rint, cint, moveType)
	if !isValid {
		return c.RespondText(errMsg)
	}
	if game.TicBoard.HasWon() {
		text := game.EndGameText() + "\n🏆برنده بازی:*" + game.Hub.Players[game.CurrentPlayerIndex].Name + "*"

		_, err := app.Bot.Edit(game.Hub, text, telebot.ModeMarkdownV2,
			CreateInlineKeyboard(
				EndgameInlineKeyboard,
			),
		)

		return err
	}
	if !game.TicBoard.IsAnyCellEmpty() {
		text := game.EndGameText() + "\nبازی مساوی شد"

		_, err := app.Bot.Edit(game.Hub, text, telebot.ModeMarkdownV2,
			CreateInlineKeyboard(
				EndgameInlineKeyboard,
			),
		)
		return err

	}

	game.NextPlayer()

	_, err := app.Bot.Edit(game.Hub, ticStartText, telebot.ModeMarkdownV2,
		CreateInlineKeyboard(
			CreateTicBoardInlineButton(game.TicBoard),
			CreatePlayersInlineButton(game.Hub.Players, game.CurrentPlayerIndex),
		),
	)

	return err

}
func (game *TicGame) EndGameText() string {
	return ticStartText + "\nبازیکن ها:\n" + game.Hub.Players[0].Name + " " + XEmoji + "\n" + game.Hub.Players[1].Name + " " + OEmoji + "\n\n" + game.CreateTicBoardAsEmoji()
}

func (game *TicGame) CreateTicBoardAsEmoji() string {
	text := ""
	for _, row := range game.TicBoard.Board {
		for _, cell := range row {
			if cell == tictactoe.Empty {
				text += EmptyEmoji
			} else if cell == tictactoe.X {
				text += XEmoji
			} else if cell == tictactoe.O {
				text += OEmoji
			}
		}
		text += "\n"
	}
	return text
}

var EndgameInlineKeyboard = [][]telebot.InlineButton{
	{
		{Text: "🤝 بازی با دوستان", InlineQueryChosenChat: &telebot.SwitchInlineQuery{AllowUserChats: true, AllowGroupChats: true}},
		{Text: "🔄 دوباره", InlineQuery: ""},
	},
}

func CreateTicBoardInlineButton(board *tictactoe.TicBoard) [][]telebot.InlineButton {
	buttons := make([][]telebot.InlineButton, 3)
	for r := range 3 {
		buttons[r] = make([]telebot.InlineButton, 3)
		for c := range 3 {
			var value string
			if board.Board[r][c] == tictactoe.Empty {
				value = "◽️"
			} else if board.Board[r][c] == tictactoe.X {
				value = "❌"
			} else {
				value = "⭕️"
			}
			buttons[r][c] = telebot.InlineButton{
				Text: value,
				Data: "xo_" + "play_" + strconv.Itoa(r) + strconv.Itoa(c),
			}
		}
	}
	return buttons
}

func CreatePlayersInlineButton(humanPlayers []*HumanPlayer, CurrentPlayerTurn int) [][]telebot.InlineButton {
	buttons := make([][]telebot.InlineButton, 0)
	for index, hplayer := range humanPlayers {

		yourTurn := ""
		if CurrentPlayerTurn == index {
			yourTurn = "🎮"
		}

		emoji := "👤"
		playEmoji := OEmoji
		if index == 0 {
			emoji = "🗿"
			playEmoji = XEmoji
		}

		name := hplayer.Name
		if len(name) > 20 {
			name = name[:20] + "..."
		}

		row := make([]telebot.InlineButton, 2)
		row = append(row, telebot.InlineButton{
			Text: fmt.Sprintf("%s %s (%s) %s", yourTurn, name, playEmoji, emoji),
			URL:  fmt.Sprintf("tg://user?id=%d", hplayer.TgID),
		})
		buttons = append(buttons, row)
	}

	return buttons
}

var startInlineKeyboard = &telebot.ReplyMarkup{
	InlineKeyboard: [][]telebot.InlineButton{
		{
			{
				Text: "واسا بازی رو بسازم",
				Data: "_",
			},
		},
	},
	ResizeKeyboard: true,
}

func (app *Application) JoinGameKeyboard() *telebot.ReplyMarkup {
	var startInlineKeyboard = &telebot.ReplyMarkup{
		InlineKeyboard: [][]telebot.InlineButton{
			{
				{
					Text: "منم بازی",
					Data: "join_hub",
				},
			},
		},
		ResizeKeyboard: true,
	}
	return startInlineKeyboard
}

const (
	EmptyEmoji = "◽️"
	XEmoji     = "❌"
	OEmoji     = "⭕️"
)
const (
	ticStartText = `❌ *دوز بازی* ⭕️`
	ticRulesText = `
	قوانین 🎮
	یک سطر یا ستون یا قطر رو با علامتت پر کن`
)
