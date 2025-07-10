package xoconsole

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/arian-nj/chibazi/database"
	consoleplayer "github.com/arian-nj/chibazi/internals/console_player"
	gametype "github.com/arian-nj/chibazi/internals/game_type"
	keybul "github.com/arian-nj/chibazi/internals/keybul"
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
	XOStartText  = `❌ *دوز بازی* ⭕️`
	ticRulesText = `
	قوانین 🎮
	یک سطر یا ستون یا قطر رو با علامتت پر کن`
)

type XOGame struct { // of GameInterface type
	XOBoard            *tictactoe.TicBoard
	CurrentPlayerIndex int
	players            []*consoleplayer.ConsolePlayer
	GameType           gametype.GameType

	ViaMessageId string // Via Bots

	Queries *database.Queries
}

func NewXOGame(gt gametype.GameType, queries *database.Queries) *XOGame {
	maxBoardSize := 3
	winSize := 3
	if gt == gametype.XOGameType5X5 {
		maxBoardSize = 5
		winSize = 4
	}
	randIndex := random.GenerateRandomNumber(2)

	return &XOGame{
		XOBoard:            tictactoe.NewTicBoard(maxBoardSize, winSize),
		players:            []*consoleplayer.ConsolePlayer{},
		CurrentPlayerIndex: randIndex,

		GameType: gt,
		Queries:  queries,
	}
}

func (g *XOGame) SendJoinPanel(c telebot.Context) error {
	sender := c.Sender()
	g.AddPlayer(consoleplayer.NewPlayer(sender.FirstName, int(sender.ID)))
	inlineKeyboard := keybul.CreateInlineKeyboard(
		keybul.JoinGameInlineButtons,
	)
	text := XOStartText + "\n\n" + g.RulesText() + "\n\n🕹 بازیکن " + fmt.Sprintf("%s", keybul.EscapeReserved(sender.FirstName)) + " منتظر حریفه"
	return g.Edit(c.Bot(), g, text, inlineKeyboard)
}

func (g *XOGame) Players() []*consoleplayer.ConsolePlayer {
	return g.players
}

func (g *XOGame) NextPlayer() {
	if g.CurrentPlayerIndex == len(g.players)-1 {
		g.CurrentPlayerIndex = 0
		return
	}
	g.CurrentPlayerIndex += 1
}

func (g *XOGame) MessageSig() (string, int64) {
	return g.ViaMessageId, 0
}

func (g *XOGame) AddPlayer(player *consoleplayer.ConsolePlayer) {
	g.players = append(g.players, player)
}

func (g *XOGame) StartGame(bot telebot.API) error {
	err := g.Edit(bot, g, XOStartText+"\n\n"+g.RulesText(),
		keybul.CreateInlineKeyboard(
			keybul.CreateBotNameInlineButton(),
			CreateTicBoardInlineButton(g.XOBoard),
			g.CreatePlayersInlineButton(g.players, g.CurrentPlayerIndex),
		),
	)
	if err != nil {
		return fmt.Errorf("error when starting xo game %w", err)
	}
	_, err = g.Queries.CreateHub(context.Background(), database.CreateHubParams{
		GameType: string(g.GameType),
		TgID:     g.players[0].TgID,
	})
	return err
}

func (g *XOGame) CallbackHandler(c telebot.Context) error {
	callbackData := c.Callback().Data
	if after, hasPrefix := strings.CutPrefix(callbackData, "join"); hasPrefix {
		return g.XOJoinGameHandler(c, after)

	} else if after, hasPrefix := strings.CutPrefix(callbackData, "play_"); hasPrefix {
		return g.XOPlayHandler(c, after)
	}
	return c.RespondAlert("no a valid callback")
}

func (g *XOGame) XOJoinGameHandler(c telebot.Context, callbackData string) error {
	sender := c.Callback().Sender
	if sender.ID == int64(g.players[0].TgID) {
		text := "خودت بازیو ساختی تو بازی هستی"
		return c.RespondText(text)
	}
	g.AddPlayer(consoleplayer.NewPlayer(sender.FirstName, int(sender.ID)))
	text := "اضافه شدی بازی شروع شد"
	err := c.RespondText(text)
	if err != nil {
		return err
	}
	return g.StartGame(c.Bot())
}

func (g *XOGame) XOPlayHandler(c telebot.Context, callbackData string) error {
	sender := c.Sender()
	if len(callbackData) != 2 {
		return fmt.Errorf("invalid ttt_play data")
	}

	if g.players[g.CurrentPlayerIndex].TgID != int(sender.ID) {
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
	if g.CurrentPlayerIndex == 0 {
		moveType = tictactoe.X
	} else {
		moveType = tictactoe.O
	}

	isValid, errMsg := g.XOBoard.PlayMove(rint, cint, moveType)
	if !isValid {
		return c.RespondText(errMsg)
	}

	hasWon := g.XOBoard.HasWon(rint, cint, moveType)
	if hasWon {
		text := g.EndGameText() + "\n🏆برنده بازی:*" + g.players[g.CurrentPlayerIndex].Name + "*"
		err := g.Edit(c.Bot(), g, text,
			keybul.CreateInlineKeyboard(
				keybul.CreateBotNameInlineButton(),
				keybul.EndGameInlineKeyboard(g.ViaMessageId != ""),
			),
		)

		return err
	}
	if !g.XOBoard.IsAnyCellEmpty() {
		text := g.EndGameText() + "\nبازی مساوی شد"

		err := g.Edit(c.Bot(), g, text,
			keybul.CreateInlineKeyboard(
				keybul.CreateBotNameInlineButton(),
				keybul.EndGameInlineKeyboard(g.ViaMessageId != ""),
			),
		)
		return err

	}
	g.NextPlayer()

	err := g.Edit(c.Bot(), g, XOStartText+"\n\n"+g.RulesText(),
		keybul.CreateInlineKeyboard(
			keybul.CreateBotNameInlineButton(),
			CreateTicBoardInlineButton(g.XOBoard),
			g.CreatePlayersInlineButton(g.players, g.CurrentPlayerIndex),
		),
	)

	return err

}

func (g *XOGame) EndGameText() string {
	return XOStartText + "\nبازیکن ها:\n" + g.players[0].Name + " " + XEmoji + "\n" + g.players[1].Name + " " + OEmoji + "\n\n" + g.CreateBoardAsEmoji()
}

func (g *XOGame) CreateBoardAsEmoji() string {
	text := ""
	for _, row := range g.XOBoard.Board {
		for _, cell := range row {
			switch cell {
			case tictactoe.Empty:
				text += EmptyEmoji
			case tictactoe.X:
				text += XEmoji
			case tictactoe.O:
				text += OEmoji
			}
		}
		text += "\n"
	}
	return text
}

func CreateTicBoardInlineButton(board *tictactoe.TicBoard) [][]telebot.InlineButton {
	buttons := make([][]telebot.InlineButton, board.MaxCellSize)
	for r := range board.MaxCellSize {
		buttons[r] = make([]telebot.InlineButton, board.MaxCellSize)
		for c := range board.MaxCellSize {
			var value string
			switch board.Board[r][c] {
			case tictactoe.Empty:
				value = "◽️"
			case tictactoe.X:
				value = "❌"
			case tictactoe.O:
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

		row := make([]telebot.InlineButton, 2)
		row = append(row, telebot.InlineButton{
			Text: fmt.Sprintf("%s %s (%s)", yourTurn, name, playEmoji),
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

func (g *XOGame) Edit(bot telebot.API, msg telebot.Editable, text string, keyboard *telebot.ReplyMarkup) error {
	if g.ViaMessageId != "" {
		err := keybul.EditGameMessage(bot, g, text, keyboard)
		if err != nil {
			return fmt.Errorf("can't edit via message %w", err)
		}
		return nil
	} else {
		for _, p := range g.players {
			err := keybul.EditGameMessage(bot, p, text, keyboard)
			if err != nil {
				slog.Error("can't edit player message ", "error", err)
			}
		}

	}
	return nil
}
