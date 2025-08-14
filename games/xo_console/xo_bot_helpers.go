package xoconsole

import (
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/arian-nj/chibazi/internals/keybul"
	"github.com/arian-nj/chibazi/internals/xo_core"
	"gopkg.in/telebot.v4"
)

func (cg *XOGame) MessageSig() (string, int64) {
	return cg.ViaMessageId, 0
}

func (g *XOGame) SendJoinPanelAddSender(c telebot.Context) error {
	sender := c.Sender()
	g.AddPlayer(sender.FirstName, int(sender.ID), nil)
	inlineKeyboard := keybul.CreateInlineKeyboard(
		keybul.JoinGameInlineButtons,
	)
	text := XOStartText + "\n\n" + g.RulesText() + "\n\n🕹 بازیکن " + fmt.Sprintf("%s", sender.FirstName) + " منتظر حریفه"
	return g.Edit(c.Bot(), g, text, inlineKeyboard)
}

const (
	XOStartText  = `❌ *دوز بازی* ⭕️`
	ticRulesText = `
	قوانین 🎮
	یک سطر یا ستون یا قطر رو با علامتت پر کن`
)
const (
	EmptyEmoji = "◽️"
	XEmoji     = "❌"  // player one
	OEmoji     = "⭕️" // player two
)

func (g *XOGame) WinGameText() string {
	return "\n🏆برنده بازی:*" + g.GetCurrentPlayer().Name + "*"

}
func (g *XOGame) EndGameText() string {
	players := g.Players
	return XOStartText + "\nبازیکن ها:\n" + players[0].Name + " " + XEmoji + "\n" + players[1].Name + " " + OEmoji + "\n\n" + g.CreateBoardAsEmoji()
}

func (g *XOGame) CreateBoardAsEmoji() string {
	text := ""
	for cellIndex, cell := range g.Board.Board {
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

func CreateTicBoardInlineButton(board *xo_core.XoBoard) [][]telebot.InlineButton {
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

func (g *XOGame) CreatePlayersInlineButton(humanPlayers []*XoPlayer, CurrentPlayerTurn int) [][]telebot.InlineButton {
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
	text += fmt.Sprintf("❕اندازه *%dX%d*\n", g.Board.MaxCellSize, g.Board.MaxCellSize)
	text += fmt.Sprintf("⚠️با یه خط *%d تایی* برنده ای", g.Board.WinSize)
	return text
}

func (g *XOGame) Edit(bot telebot.API, msg telebot.Editable, text string, keyboard *telebot.ReplyMarkup) error {
	g.LastEdit = time.Now()
	if g.ViaMessageId != "" {
		err := keybul.EditGameMessage(bot, g, text, keyboard)
		if err != nil {
			return fmt.Errorf("can't edit via message %w", err)
		}
		return nil
	} else {
		for _, p := range g.Players {
			if p.MessageID == 0 {
				msg, err := bot.Send(p, "game")
				if err != nil {
					slog.Error("can't send player message ", "error", err)
				}
				p.SetMessageSig(msg.ID)
			}
			err := keybul.EditGameMessage(bot, p, text, keyboard)
			if err != nil {
				slog.Error("can't edit player message ", "error", err)
			}
		}

	}
	return nil
}
