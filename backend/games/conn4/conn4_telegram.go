package conn4

import (
	"fmt"
	"log/slog"
	"slices"
	"strconv"
	"strings"
	"time"

	conn4_core "github.com/arian-nj/chibazi/backend/games/conn4/core"
	"github.com/arian-nj/chibazi/backend/internals/commander"
	"github.com/arian-nj/chibazi/backend/internals/keybul"
	"gopkg.in/telebot.v4"
)

type Conn4TelegramListener struct {
	Bot *telebot.Bot

	LastEdit     time.Time
	Player       *Conn4Player
	ViaMessageId string // Via Bots
}

func newConn4TelegramListener(player *Conn4Player, bot *telebot.Bot, viaMessageID string) *Conn4TelegramListener {
	return &Conn4TelegramListener{
		Bot:          bot,
		Player:       player,
		LastEdit:     time.Now(),
		ViaMessageId: viaMessageID,
	}
}
func (game *Conn4State) CallBackRouter(c telebot.Context) error {
	callbackData := c.Callback().Data
	if after, hasPrefix := strings.CutPrefix(callbackData, "play_"); hasPrefix {
		return game.Conn4PlayHandler(c, after)
	}
	slog.Error("invalid callback in xo game callback router")
	return c.RespondAlert("no a valid callback")
}

func (tg *Conn4TelegramListener) Update(command commander.Command) {
	switch c := command.(type) {
	case *PlayCommand:
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
			err := tg.TieGame(c)
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

func (game *Conn4State) Conn4PlayHandler(c telebot.Context, callbackData string) error {
	sender := c.Sender()
	currentPlayer := game.CurrentPlayer()

	if currentPlayer.TelegramID != int(sender.ID) {
		return c.RespondText("نوبت تو نیست!")
	}

	rowIndex, err := strconv.Atoi(callbackData)
	if err != nil {
		return c.RespondAlert("یه مشکلی هست")
	}

	moveType := currentPlayer.Move

	isValid, errMsg := game.Board.IsMoveValid(rowIndex)
	if !isValid {
		return c.RespondText(errMsg)
	}

	player := game.findByTelegramID(int(sender.ID))
	if player == nil {
		return c.RespondText("can't find player")
	}

	playCommand := NewPlayCommand(game, rowIndex, moveType, player.ID)
	game.PushCommand(playCommand)

	return nil
}

// send Commands
var Conn4PlayNumberButton = [][]telebot.InlineButton{
	{
		{Text: "1", Data: "play_0"},
		{Text: "2", Data: "play_1"},
		{Text: "3", Data: "play_2"},
		{Text: "4", Data: "play_3"},
		{Text: "5", Data: "play_4"},
		{Text: "6", Data: "play_5"},
		{Text: "7", Data: "play_6"},
	},
}

func MakeBoardAsEmojies(game *Conn4State, winList []int) string {
	board := game.Board.Board
	boardText := ""
	winListLen := len(winList)
	for cIndex, cell := range board {
		cellEmoji := EmptyEmoji
		switch cell {
		case conn4_core.Empty:
		case conn4_core.One:
			cellEmoji = OneEmoji
		case conn4_core.Two:
			cellEmoji = TwoEmoji
		}
		if winListLen != 0 && slices.Contains(winList, cIndex) {
			cellEmoji = WinEmoji
		}
		boardText += cellEmoji
		if cIndex%(conn4_core.BOARD_WIDTH) == conn4_core.BOARD_WIDTH-1 {
			boardText += "\n"
		}
	}
	return boardText
}

func (tg *Conn4TelegramListener) EditDuringGameBoard(game *Conn4State) error {
	err := tg.Edit(tg, Conn4StartText+"\n"+MakeBoardAsEmojies(game, nil)+
		"1️⃣2️⃣3️⃣4️⃣5️⃣6️⃣7️⃣",
		keybul.CreateInlineKeyboard(
			// keybul.ContinueInWebButton(),
			Conn4PlayNumberButton,
			game.CreatePlayersInlineButton(game.Players, game.CurrentPlayerIndex),
		),
		tg.Player,
	)
	return err
}

func (tg *Conn4TelegramListener) Edit(msg telebot.Editable, text string, keyboard *telebot.ReplyMarkup, player *Conn4Player) error {
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

func (tg *Conn4TelegramListener) MessageSig() (string, int64) {
	return tg.ViaMessageId, 0
}

const (
	OneEmoji   = "🔵"
	TwoEmoji   = "🔴"
	EmptyEmoji = "⚪️"
	WinEmoji   = "🟢"
)

func (game *Conn4State) CreatePlayersInlineButton(humanPlayers []*Conn4Player, CurrentPlayerTurn int) [][]telebot.InlineButton {
	rows := make([]telebot.InlineButton, 2)
	for index, hplayer := range humanPlayers {

		yourTurn := ""
		if CurrentPlayerTurn == index {
			yourTurn = "🎮"
		}

		// emoji := ""
		playEmoji := OneEmoji
		if index == 1 {
			// emoji = "🗿"
			playEmoji = TwoEmoji
		}

		name := hplayer.Name
		if len(name) > 20 {
			name = name[:20] + "..."
		}

		remainedTime := MAX_ALLOWED_TIME - hplayer.Timer.Spent()
		timeText := fmt.Sprintf("%d:%d", int(remainedTime.Minutes()), int(remainedTime.Seconds())%60)

		btn := telebot.InlineButton{
			Text: fmt.Sprintf("%s %s (%s) %s", yourTurn, name, playEmoji, timeText),
			Data: "_",
		}
		rows = append(rows, btn)
	}

	return [][]telebot.InlineButton{
		rows,
	}
}

func WinnerGameText(winner *Conn4Player) string {
	return "\n🏆برنده بازی:" + "*" + winner.Name + "*"
}
func (game *Conn4State) EndGameText(winLine []int) string {
	players := game.Players
	return Conn4StartText + "\nبازیکن ها:\n" +
		players[0].Name + " " + OneEmoji + "\n" + players[1].Name + " " + TwoEmoji + "\n\n" +
		MakeBoardAsEmojies(game, winLine)
}

func (tg *Conn4TelegramListener) TheEnd(endCommand *EndGameCommand) error {
	game := endCommand.Game
	additionalText := ""
	if endCommand.reason == END_GAME_TIMEOUT {
		additionalText = "\n برنده زمانی"
	}

	text := game.EndGameText(endCommand.WinLine) + WinnerGameText(endCommand.Winner) + additionalText
	err := tg.Edit(tg, text,
		keybul.CreateInlineKeyboard(
			// keybul.ContinueInWebButton(),
			keybul.EndGameInlineKeyboard(tg.ViaMessageId != ""),
		),
		tg.Player,
	)
	return err
}

func (tg *Conn4TelegramListener) TieGame(endCommand *EndGameCommand) error {
	text := endCommand.Game.EndGameText(endCommand.WinLine) + "\nبازی مساوی شد"

	err := tg.Edit(tg, text,
		keybul.CreateInlineKeyboard(
			keybul.ContinueInWebButton(),
			keybul.EndGameInlineKeyboard(tg.ViaMessageId != ""),
		),
		tg.Player,
	)
	return err
}
