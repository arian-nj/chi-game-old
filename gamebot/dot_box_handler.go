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
	*Game
	DotBoxBoard    *dotbox.DotBoxBoard
	PlayerOnScore  int
	PlayerTwoScore int
}

func NewDotBoxGame() *DotBoxGame {
	return &DotBoxGame{
		DotBoxBoard: dotbox.NewDotBoard(),
		Game:        NewGame(),
	}
}

const (
	dotBoxStartText = `❌ *نقطه بازی* ⭕️`
)

func (game *DotBoxGame) StartGame(bot *telebot.Bot) {
	randIndex := random.GenerateRandomNumber(2)
	game.CurrentPlayerIndex = randIndex

	_, err := bot.Edit(game, ticStartText, telebot.ModeMarkdownV2,
		CreateInlineKeyboard(
			CreateBotNameInlineButton(bot),
			CreateDotBoxInlineButton(game.DotBoxBoard),
			game.CreatePlayersInlineButton(game.Players, game.CurrentPlayerIndex),
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

func (game *DotBoxGame) CallbackHandlers(c telebot.Context, bot *telebot.Bot) error {
	callback := c.Callback()
	callbackData := strings.TrimPrefix(callback.Data, string(DotBoxGameType)+"_")
	if strings.HasPrefix(callbackData, "play_") {
		return game.PlayHandler(c, bot)
	}
	return nil
}

func (game *DotBoxGame) PlayHandler(c telebot.Context, bot *telebot.Bot) error {
	callback := c.Callback()
	callbackData := strings.TrimPrefix(callback.Data, string(DotBoxGameType)+"_play_")

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

	if !game.DotBoxBoard.HasEmptyCell() {
		text := ""
		if game.PlayerOnScore != game.PlayerTwoScore {
			winnerPlayer := game.Players[0]
			if game.PlayerTwoScore > game.PlayerOnScore {
				winnerPlayer = game.Players[0]
			}
			text = "\n🏆برنده بازی:*" + winnerPlayer.Name + "*"
		} else {
			text = "\nبازی مساوی شد"

		}

		_, err := bot.Edit(game, game.EndGameText()+text, telebot.ModeMarkdownV2,
			CreateInlineKeyboard(
				CreateBotNameInlineButton(bot),
				EndgameInlineKeyboard,
			),
		)
		return err

	}
	_, err := bot.Edit(game, dotBoxStartText+"\n"+game.CreateBoardAsEmoji(), telebot.ModeMarkdownV2,
		CreateInlineKeyboard(
			CreateBotNameInlineButton(bot),
			CreateDotBoxInlineButton(game.DotBoxBoard),
			game.CreatePlayersInlineButton(game.Players, game.CurrentPlayerIndex),
		),
	)

	if err != nil {
		return err
	}
	return nil
}

func (game *DotBoxGame) EndGameText() string {
	return ticStartText + "\nبازیکن ها:\n" + game.Players[0].Name + " " + strconv.Itoa(game.PlayerOnScore) + "\n" + game.Players[1].Name + " " + strconv.Itoa(game.PlayerTwoScore)
}

func (game *DotBoxGame) CreateBoardAsEmoji() string {
	text := ""
	for r, row := range game.DotBoxBoard.Board {
		for c, cell := range row {
			value := "⌾"
			if cell != dotbox.Empty {
				value = getSymbol(r, c)
			}
			text += value
		}
		text += "\n"
	}
	return text
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
