package xoconsole

import (
	"fmt"
	"strconv"

	consoleplayer "github.com/arian-nj/chibazi/internals/console_player"
	"github.com/arian-nj/chibazi/internals/xo_core"
	"gopkg.in/telebot.v4"
)

func (g *XOGame) WinGameText() string {
	return "\n🏆برنده بازی:*" + g.GetCurrentPlayer().Name + "*"

}
func (g *XOGame) EndGameText() string {
	players := g.Players()
	return XOStartText + "\nبازیکن ها:\n" + players[0].Name + " " + XEmoji + "\n" + players[1].Name + " " + OEmoji + "\n\n" + g.CreateBoardAsEmoji()
}

func (g *XOGame) CreateBoardAsEmoji() string {
	text := ""
	for cellIndex, cell := range g.XOBoard.Board {
		if cellIndex%3 == 0 {
			text += "\n"
		}

		switch cell {
		case xo_core.Empty:
			text += EmptyEmoji
		case xo_core.X:
			text += XEmoji
		case xo_core.O:
			text += OEmoji
		}
	}

	return text
}

func CreateTicBoardInlineButton(board *xo_core.TicBoard) [][]telebot.InlineButton {
	buttons := make([][]telebot.InlineButton, board.MaxCellSize)
	for r := range board.MaxCellSize {
		buttons[r] = make([]telebot.InlineButton, board.MaxCellSize)
		for c := range board.MaxCellSize {
			var value string
			switch board.GetCell(r, c) {
			case xo_core.Empty:
				value = "◽️"
			case xo_core.X:
				value = "❌"
			case xo_core.O:
				value = "⭕️"
			}

			buttons[r][c] = telebot.InlineButton{
				Text: value,
				Data: "play_" + strconv.Itoa(r) + strconv.Itoa(c),
			}
		}
	}
	return buttons
}

func (g *XOGame) CreatePlayersInlineButton(humanPlayers []*consoleplayer.ConsolePlayer, CurrentPlayerTurn int) [][]telebot.InlineButton {
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

		remainedTime := MaxPlayerTime - hplayer.SpentTime
		timeText := fmt.Sprintf("%d:%d", int(remainedTime.Minutes()), int(remainedTime.Seconds())%60)

		row := make([]telebot.InlineButton, 2)
		row = append(row, telebot.InlineButton{
			Text: fmt.Sprintf("%s %s (%s) %s", yourTurn, name, playEmoji, timeText),
			// URL:  fmt.Sprintf("tg://user?id=%d", hplayer.TgID),
			Data: "_",
		})
		buttons = append(buttons, row)
	}

	return buttons
}

func (g *XOGame) RulesText() string {
	text := ""
	// text += "قوانین:\د"
	text += fmt.Sprintf("❕اندازه *%dX%d*\n", g.XOBoard.MaxCellSize, g.XOBoard.MaxCellSize)
	text += fmt.Sprintf("⚠️با یه خط *%d تایی* برنده ای", g.XOBoard.WinSize)
	return text
}
