package xo

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/arian-nj/chibazi/backend/internals/commander"
	"github.com/arian-nj/chibazi/backend/internals/keybul"
	"gopkg.in/telebot.v4"
)

type XoTelegramListener struct {
	Bot *telebot.Bot

	LastEdit     time.Time
	Player       *XoPlayer
	ViaMessageId string // Via Bots
}

func newXOTelegramListener(player *XoPlayer, bot *telebot.Bot, viaMessageID string) *XoTelegramListener {
	return &XoTelegramListener{
		Bot:          bot,
		Player:       player,
		LastEdit:     time.Now(),
		ViaMessageId: viaMessageID,
	}
}

func (tg *XoTelegramListener) Update(command commander.Command) {
	switch c := command.(type) {
	case *MoveCommand:
		err := tg.EditDuringGameBoard(c.Game)
		if err != nil {
			slog.Error("can't edit durning move command", "err", err)
		}
	case *StartCommand:
		err := tg.EditDuringGameBoard(c.Game)
		if err != nil {
			slog.Error("can't edit during start command", "err", err)
		}
	case *EndGameCommand:
		if c.Winner == nil {
			err := tg.TieGame(c.Game)
			if err != nil {
				slog.Error("tie failed", "error", err)
			}
		} else {
			err := tg.TheEnd(c)
			if err != nil {
				slog.Error("the end failed", "error", err)
			}
		}
	case *SyncTimeCommand:
		if time.Since(tg.LastEdit) > 10*time.Second {
			err := tg.EditDuringGameBoard(c.Game)
			if err != nil {
				slog.Error("can't edit sync Time", "err", err)
			}
		}
	}
}

func (g *XOState) CallBackRouter(c telebot.Context) error {
	callbackData := c.Callback().Data
	if after, hasPrefix := strings.CutPrefix(callbackData, "play_"); hasPrefix {
		return g.XOPlayHandler(c, after)
	}
	slog.Error("invalid callback in xo game callback router")
	return c.RespondAlert("no a valid callback")
}

// send Commands
func (tg *XoTelegramListener) EditDuringGameBoard(game *XOState) error {
	err := tg.Edit(tg, XOStartText+"\n\n"+game.RulesText(),
		keybul.CreateInlineKeyboard(
			keybul.ContinueInWebButton(),
			CreateTicBoardInlineButton(game.Board),
			game.CreatePlayersInlineButton(game.Players, game.CurrentPlayerIndex),
		),
		tg.Player,
	)
	return err
}

func (tg *XoTelegramListener) TheEnd(endCommand *EndGameCommand) error {
	game := endCommand.Game
	additionalText := ""
	if endCommand.reason == END_GAME_TIMEOUT {
		additionalText = "\n برنده زمانی"
	}

	text := game.EndGameText() + game.WinGameText(endCommand.Winner) + additionalText
	err := tg.Edit(tg, text,
		keybul.CreateInlineKeyboard(
			keybul.ContinueInWebButton(),
			keybul.EndGameInlineKeyboard(tg.ViaMessageId != ""),
		),
		tg.Player,
	)
	return err
}

func (tg *XoTelegramListener) TieGame(game *XOState) error {
	text := game.EndGameText() + "\nبازی مساوی شد"

	err := tg.Edit(tg, text,
		keybul.CreateInlineKeyboard(
			keybul.ContinueInWebButton(),
			keybul.EndGameInlineKeyboard(tg.ViaMessageId != ""),
		),
		tg.Player,
	)
	return err
}

// Reciecves Actions

func (game *XOState) XOPlayHandler(c telebot.Context, callbackData string) error {
	sender := c.Sender()

	if game.CurrentPlayer().TelegramID != int(sender.ID) {
		return c.RespondText("نوبت تو نیست!")
	}

	cellIndex, err := strconv.Atoi(callbackData)
	if err != nil {
		c.RespondAlert("یه مشکلی هست")
	}

	moveType := game.CurrentPlayer().Move

	isValid, errMsg := game.Board.IsMoveValid(cellIndex, moveType)
	if !isValid {
		return c.RespondText(errMsg)
	}

	player := game.findByTelegramID(int(sender.ID))
	if player == nil {
		return c.RespondText("can't find player")
	}

	playCommand := NewPlayCommand(game, cellIndex, moveType, player.ID)
	game.PushCommand(playCommand)
	return nil
}

// HELPERS

const (
	EmptyEmoji = "◽️"
	XEmoji     = "❌"  // player one
	OEmoji     = "⭕️" // player two
)

func (game *XOState) WinGameText(winner *XoPlayer) string {
	return "\n🏆برنده بازی:" + "*" + winner.Name + "*"

}
func (game *XOState) EndGameText() string {
	players := game.Players
	return XOStartText + "\nبازیکن ها:\n" + players[0].Name + " " + XEmoji + "\n" + players[1].Name + " " + OEmoji + "\n\n" + game.CreateBoardAsEmoji()
}

var xoBoardSymbols = map[Cell]string{
	Empty: "◽️",
	X:     "❌",
	O:     "⭕️",
}

func (game *XOState) CreateBoardAsEmoji() string {
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

func (game *XOState) CreatePlayersInlineButton(humanPlayers []*XoPlayer, CurrentPlayerTurn int) [][]telebot.InlineButton {
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

		remainedTime := MAX_ALLOWED_TIME - hplayer.Timer.Spent()
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

func (tg *XoTelegramListener) Edit(msg telebot.Editable, text string, keyboard *telebot.ReplyMarkup, player *XoPlayer) error {
	tg.LastEdit = time.Now()
	if tg.ViaMessageId != "" {
		err := keybul.EditMessage(tg.Bot, tg, text, keyboard)
		if err != nil {
			return fmt.Errorf("can't edit via message %w", err)
		}
		return nil
	} else {
		if player.MessageID == 0 {
			msg, err := tg.Bot.Send(player, "game")
			if err != nil {
				slog.Error("can't send player message ", "error", err)
			}
			player.MessageID = msg.ID
		}
		err := keybul.EditMessage(tg.Bot, player, text, keyboard)
		if err != nil {
			slog.Error("can't edit player message ", "error", err)
		}

	}
	return nil
}

func (tg *XoTelegramListener) MessageSig() (string, int64) {
	return tg.ViaMessageId, 0
}
