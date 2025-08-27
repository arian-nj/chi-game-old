package xo

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/arian-nj/chibazi/backend/internals/keybul"
	"gopkg.in/telebot.v4"
)

type XoTelegramListener struct {
	Bot *telebot.Bot

	LastEdit     time.Time
	ViaMessageId string // Via Bots
}

func NewXOTelegramListener(bot *telebot.Bot, viaMessageID string) *XoTelegramListener {
	return &XoTelegramListener{
		Bot:          bot,
		LastEdit:     time.Now(),
		ViaMessageId: viaMessageID,
	}
}

func (tg *XoTelegramListener) Update(game *XOGame, command Command) {
	return
	switch c := command.(type) {
	case *MoveCommand:
		tg.EditDuringGameBoard(game)
	case *StartCommand:
		tg.EditDuringGameBoard(game)
	case *EndGameCommand:
		if c.Winner == nil {
			err := tg.TieGame(game)
			if err != nil {
				slog.Error("tie failed", "error", err)
			}
		} else {
			err := tg.TheEnd(game, c.Winner, c.Text)
			if err != nil {
				slog.Error("the end failed", "error", err)
			}
		}
	}
}

func (g *XOGame) CallBackRouter(c telebot.Context) error {
	callbackData := c.Callback().Data
	if callbackData == "join" {
		// return g.XOJoinGameHandler(c)

	} else if after, hasPrefix := strings.CutPrefix(callbackData, "play_"); hasPrefix {
		return g.XOPlayHandler(c, after)
	}
	return c.RespondAlert("no a valid callback")
}

// send Commands
func (tg *XoTelegramListener) EditDuringGameBoard(game *XOGame) error {
	err := tg.Edit(tg, XOStartText+"\n\n"+game.RulesText(),
		keybul.CreateInlineKeyboard(
			keybul.CreateBotNameInlineButton(),
			CreateTicBoardInlineButton(game.Board),
			game.CreatePlayersInlineButton(game.Players, game.CurrentPlayerIndex),
		),
		game.Players,
	)
	return err
}

func (tg *XoTelegramListener) TheEnd(game *XOGame, winner *XoPlayer, additionalText string) error {
	text := game.EndGameText() + game.WinGameText(winner) + additionalText
	err := tg.Edit(tg, text,
		keybul.CreateInlineKeyboard(
			keybul.CreateBotNameInlineButton(),
			keybul.EndGameInlineKeyboard(tg.ViaMessageId != ""),
		),
		game.Players,
	)
	return err
}

func (tg *XoTelegramListener) TieGame(game *XOGame) error {
	text := game.EndGameText() + "\nبازی مساوی شد"

	err := tg.Edit(tg, text,
		keybul.CreateInlineKeyboard(
			keybul.CreateBotNameInlineButton(),
			keybul.EndGameInlineKeyboard(tg.ViaMessageId != ""),
		),
		game.Players,
	)
	return err
}

// Reciecves Actions

// func (g *XOGame) XOJoinGameHandler(c telebot.Context) error {
// 	sender := c.Callback().Sender
// 	if sender.ID == int64(g.Players[0].TgID) {
// 		text := "خودت بازیو ساختی تو بازی هستی"
// 		return c.RespondText(text)
// 	}
// 	g.AddPlayer(sender.FirstName, int(sender.ID), nil)
// 	text := "اضافه شدی بازی شروع شد"
// 	err := c.RespondText(text)
// 	if err != nil {
// 		return err
// 	}
// 	return g.StartGame()
// }

func (game *XOGame) XOPlayHandler(c telebot.Context, callbackData string) error {
	sender := c.Sender()

	if game.getCurrentPlayer().TelegramID != int(sender.ID) {
		return c.RespondText("نوبت تو نیست!")
	}

	cellIndex, err := strconv.Atoi(callbackData)
	if err != nil {
		c.RespondAlert("یه مشکلی هست")
	}

	moveType := game.getCurrentPlayer().Move

	isValid, errMsg := game.Board.IsMoveValid(cellIndex, moveType)
	if !isValid {
		return c.RespondText(errMsg)
	}

	player := game.findByTelegramID(int(sender.ID))
	if player == nil {
		return c.RespondText("can't find player")
	}

	playCommand := NewPlayCommand(cellIndex, moveType, player.ID)
	game.pushCommand(playCommand)
	return nil
}

// HELPERS

func (tg *XoTelegramListener) MessageSig() (string, int64) {
	return tg.ViaMessageId, 0
}

func (game *XOGame) SendJoinPanelAddSender(c telebot.Context) error {
	// sender := c.Sender()
	// game.AddPlayer(sender.FirstName, int(sender.ID), nil)
	// inlineKeyboard := keybul.CreateInlineKeyboard(
	// 	keybul.JoinGameInlineButtons,
	// )
	// text := XOStartText + "\n\n" + game.RulesText() + "\n\n🕹 بازیکن " + fmt.Sprintf("%s", sender.FirstName) + " منتظر حریفه"
	// return game.Edit(c.Bot(), game, text, inlineKeyboard)
	return nil
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

func (game *XOGame) WinGameText(winner *XoPlayer) string {
	return "\n🏆برنده بازی:" + "*" + winner.Name + "*"

}
func (game *XOGame) EndGameText() string {
	players := game.Players
	return XOStartText + "\nبازیکن ها:\n" + players[0].Name + " " + XEmoji + "\n" + players[1].Name + " " + OEmoji + "\n\n" + game.CreateBoardAsEmoji()
}

var xoBoardSymbols = map[Cell]string{
	Empty: "◽️",
	X:     "❌",
	O:     "⭕️",
}

func (game *XOGame) CreateBoardAsEmoji() string {
	text := ""
	for cellIndex, cell := range game.Board.Board {
		if cellIndex%3 == 0 {
			text += "\n"
		}

		text += xoBoardSymbols[cell]
	}
	return text
}

func CreateTicBoardInlineButton(board *XoBoard) [][]telebot.InlineButton {
	buttons := make([][]telebot.InlineButton, board.MaxCellSize)
	for r := range board.MaxCellSize {
		buttons[r] = make([]telebot.InlineButton, board.MaxCellSize)
		for c := range board.MaxCellSize {
			cellIndex := board.toCellIndex(r, c)
			cell := board.GetCell(cellIndex)

			buttons[r][c] = telebot.InlineButton{
				Text: xoBoardSymbols[cell],
				Data: "play_" + strconv.Itoa(cellIndex),
			}
		}
	}
	return buttons
}

func (game *XOGame) CreatePlayersInlineButton(humanPlayers []*XoPlayer, CurrentPlayerTurn int) [][]telebot.InlineButton {
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

func (game *XOGame) RulesText() string {
	text := ""
	// text += "قوانین:\د"
	text += fmt.Sprintf("❕اندازه *%dX%d*\n", game.Board.MaxCellSize, game.Board.MaxCellSize)
	text += fmt.Sprintf("⚠️با یه خط *%d تایی* برنده ای", game.Board.WinSize)
	return text
}

func (tg *XoTelegramListener) Edit(msg telebot.Editable, text string, keyboard *telebot.ReplyMarkup, players []*XoPlayer) error {
	tg.LastEdit = time.Now()
	if tg.ViaMessageId != "" {
		err := keybul.EditGameMessage(tg.Bot, tg, text, keyboard)
		if err != nil {
			return fmt.Errorf("can't edit via message %w", err)
		}
		return nil
	} else {
		for _, p := range players {
			if p.MessageID == 0 {
				msg, err := tg.Bot.Send(p, "game")
				if err != nil {
					slog.Error("can't send player message ", "error", err)
				}
				p.MessageID = msg.ID
			}
			err := keybul.EditGameMessage(tg.Bot, p, text, keyboard)
			if err != nil {
				slog.Error("can't edit player message ", "error", err)
			}
		}

	}
	return nil
}
