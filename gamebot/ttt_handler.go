package gamebot

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/arian-nj/chibazi/internals/random"
	tictactoe "github.com/arian-nj/chibazi/internals/tic_tac_toe"
	"gopkg.in/telebot.v4"
)

const (
	EmptyEmoji = "◽️"
	XEmoji     = "❌"  // player one
	OEmoji     = "⭕️" // player two
)

const (
	ticStartText = `❌ *دوز بازی* ⭕️`
	ticRulesText = `
	قوانین 🎮
	یک سطر یا ستون یا قطر رو با علامتت پر کن`
)

type TicGame struct { // of GameInterface type
	*Game
	TicBoard *tictactoe.TicBoard
}

func NewTicGame() *TicGame {
	return &TicGame{
		TicBoard: tictactoe.NewTicBoard(),
		Game:     NewGame(),
	}
}

func (game *TicGame) StartGame(bot *telebot.Bot) {
	randIndex := random.GenerateRandomNumber(2)
	game.CurrentPlayerIndex = randIndex

	_, err := bot.Edit(game, ticStartText, telebot.ModeMarkdownV2,
		CreateInlineKeyboard(
			CreateBotNameInlineButton(bot),
			CreateTicBoardInlineButton(game.TicBoard),
			game.CreatePlayersInlineButton(game.Players, game.CurrentPlayerIndex),
		),
	)
	if err != nil {
		slog.Error("error when starting ttt game" + err.Error())
		return
	}
}

func (game *TicGame) GetGameType() GameType {
	return TicTacToeGameType
}

func (game *TicGame) GetMaxPlayer() int {
	return 2
}

func (game *TicGame) TTTCallbackHandlers(c telebot.Context, bot *telebot.Bot) error {
	callback := c.Callback()
	callbackData := strings.TrimPrefix(callback.Data, string(TicTacToeGameType)+"_")
	if strings.HasPrefix(callbackData, "play_") {
		return game.TicPlayHandler(c, bot)
	}
	return nil
}

func (game *TicGame) TicPlayHandler(c telebot.Context, bot *telebot.Bot) error {
	callback := c.Callback()
	callbackData := strings.TrimPrefix(callback.Data, string(TicTacToeGameType)+"_play_")

	if len(callbackData) != 2 {
		return fmt.Errorf("invalid ttt_play data")
	}

	if game.Players[game.CurrentPlayerIndex].TgID != int(callback.Sender.ID) {
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
		text := game.EndGameText() + "\n🏆برنده بازی:*" + game.Players[game.CurrentPlayerIndex].Name + "*"

		_, err := bot.Edit(game, text, telebot.ModeMarkdownV2,
			CreateInlineKeyboard(
				CreateBotNameInlineButton(bot),
				EndgameInlineKeyboard,
			),
		)

		return err
	}
	if !game.TicBoard.IsAnyCellEmpty() {
		text := game.EndGameText() + "\nبازی مساوی شد"

		_, err := bot.Edit(game, text, telebot.ModeMarkdownV2,
			CreateInlineKeyboard(
				CreateBotNameInlineButton(bot),
				EndgameInlineKeyboard,
			),
		)
		return err

	}
	game.NextPlayer()

	_, err := bot.Edit(game, ticStartText, telebot.ModeMarkdownV2,
		CreateInlineKeyboard(
			CreateBotNameInlineButton(bot),
			CreateTicBoardInlineButton(game.TicBoard),
			game.CreatePlayersInlineButton(game.Players, game.CurrentPlayerIndex),
		),
	)

	return err

}

func (game *TicGame) EndGameText() string {
	return ticStartText + "\nبازیکن ها:\n" + game.Players[0].Name + " " + XEmoji + "\n" + game.Players[1].Name + " " + OEmoji + "\n\n" + game.CreateBoardAsEmoji()
}

func (game *TicGame) CreateBoardAsEmoji() string {
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
				Data: string(TicTacToeGameType) + "_play_" + strconv.Itoa(r) + strconv.Itoa(c),
			}
		}
	}
	return buttons
}

func (game *TicGame) CreatePlayersInlineButton(humanPlayers []*HumanPlayer, CurrentPlayerTurn int) [][]telebot.InlineButton {
	buttons := make([][]telebot.InlineButton, 0)
	for index, hplayer := range humanPlayers {

		yourTurn := ""
		if CurrentPlayerTurn == index {
			yourTurn = "🎮"
		}

		// emoji := ""
		playEmoji := OEmoji
		if index == 0 {
			// emoji = "🗿"
			playEmoji = XEmoji
		}

		name := hplayer.Name
		if len(name) > 20 {
			name = name[:20] + "..."
		}

		row := make([]telebot.InlineButton, 2)
		row = append(row, telebot.InlineButton{
			Text: fmt.Sprintf("%s %s (%s)", yourTurn, name, playEmoji),
			URL:  fmt.Sprintf("tg://user?id=%d", hplayer.TgID),
		})
		buttons = append(buttons, row)
	}

	return buttons
}
