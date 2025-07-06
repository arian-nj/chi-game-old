package xoconsole

import (
	"fmt"
	"strconv"
	"strings"
	"time"

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

type Player struct {
	TgID int
	Name string
}

func newPlayer(name string, tgID int) *Player {
	return &Player{
		Name: name,
		TgID: tgID,
	}
}

type XOGame struct { // of GameInterface type
	TicBoard           *tictactoe.TicBoard
	CurrentPlayerIndex int
	Players            []*Player
	MessageId          string
	GameType           gametype.GameType

	CreatedAt time.Time
}

func NewXOGame(gt gametype.GameType) *XOGame {
	maxBoardSize := 3
	winSize := 3
	if gt == gametype.XOGameType5X5 {
		maxBoardSize = 5
		winSize = 4
	}

	return &XOGame{
		TicBoard: tictactoe.NewTicBoard(maxBoardSize, winSize),
		Players:  []*Player{},

		GameType: gt,

		CreatedAt: time.Now(),
	}
}

func (g *XOGame) SendJoinPanel(c telebot.Context) error {
	sender := c.Sender()
	g.addPlayer(newPlayer(sender.FirstName, int(sender.ID)))
	inlineKeyboard := keybul.CreateInlineKeyboard(
		keybul.JoinGameKeyboard(g.GameType),
	)
	text := XOStartText + "\n\n🕹 بازیکن " + fmt.Sprintf("[%s](tg://user?id=%d)", sender.FirstName, sender.ID) + " منتظر حریفه"
	return keybul.EditGameMessage(c.Bot(), g, text, inlineKeyboard)
}

func (g *XOGame) NextPlayer() {
	if g.CurrentPlayerIndex == len(g.Players)-1 {
		g.CurrentPlayerIndex = 0
		return
	}
	g.CurrentPlayerIndex += 1
}

func (g *XOGame) MessageSig() (string, int64) {
	return g.MessageId, 0
}

func (g *XOGame) addPlayer(player *Player) {
	g.Players = append(g.Players, player)
}

func (game *XOGame) StartGame(c telebot.Context) error {
	randIndex := random.GenerateRandomNumber(2)
	game.CurrentPlayerIndex = randIndex

	_, err := c.Bot().Edit(game, XOStartText, telebot.ModeMarkdownV2,
		keybul.CreateInlineKeyboard(
			keybul.CreateBotNameInlineButton(),
			CreateTicBoardInlineButton(game.TicBoard),
			game.CreatePlayersInlineButton(game.Players, game.CurrentPlayerIndex),
		),
	)
	if err != nil {
		return fmt.Errorf("error when starting xo game %w", err)
	}
	return nil
}

func (game *XOGame) XOCallbackHandlers(c telebot.Context, callbackData string) error {
	if after, hasPrefix := strings.CutPrefix(callbackData, "join"); hasPrefix {
		return game.XOJoinGameHandler(c, after)

	} else if after, hasPrefix := strings.CutPrefix(callbackData, "play_"); hasPrefix {
		return game.XOPlayHandler(c, after)
	}
	return c.RespondAlert("no a valid callback")
}

func (game *XOGame) XOJoinGameHandler(c telebot.Context, callbackData string) error {
	sender := c.Callback().Sender
	if sender.ID == int64(game.Players[0].TgID) {
		text := "خودت بازیو ساختی تو بازی هستی"
		return c.RespondText(text)
	}
	game.addPlayer(newPlayer(sender.FirstName, int(sender.ID)))
	text := "اضافه شدی بازی شروع شد"
	err := c.RespondText(text)
	if err != nil {
		return err
	}
	return game.StartGame(c)
}

func (game *XOGame) XOPlayHandler(c telebot.Context, callbackData string) error {
	sender := c.Sender()
	if len(callbackData) != 2 {
		return fmt.Errorf("invalid ttt_play data")
	}

	if game.Players[game.CurrentPlayerIndex].TgID != int(sender.ID) {
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
	hasWon := game.TicBoard.HasWon(rint, cint, moveType)
	if hasWon {
		text := game.EndGameText() + "\n🏆برنده بازی:*" + game.Players[game.CurrentPlayerIndex].Name + "*"
		_, err := c.Bot().Edit(game, text, telebot.ModeMarkdownV2,
			keybul.CreateInlineKeyboard(
				keybul.CreateBotNameInlineButton(),
				keybul.EndgameInlineKeyboard,
			),
		)

		return err
	}
	if !game.TicBoard.IsAnyCellEmpty() {
		text := game.EndGameText() + "\nبازی مساوی شد"

		_, err := c.Bot().Edit(game, text, telebot.ModeMarkdownV2,
			keybul.CreateInlineKeyboard(
				keybul.CreateBotNameInlineButton(),
				keybul.EndgameInlineKeyboard,
			),
		)
		return err

	}
	game.NextPlayer()

	_, err := c.Bot().Edit(game, XOStartText, telebot.ModeMarkdownV2,
		keybul.CreateInlineKeyboard(
			keybul.CreateBotNameInlineButton(),
			CreateTicBoardInlineButton(game.TicBoard),
			game.CreatePlayersInlineButton(game.Players, game.CurrentPlayerIndex),
		),
	)

	return err

}

func (game *XOGame) EndGameText() string {
	return XOStartText + "\nبازیکن ها:\n" + game.Players[0].Name + " " + XEmoji + "\n" + game.Players[1].Name + " " + OEmoji + "\n\n" + game.CreateBoardAsEmoji()
}

func (game *XOGame) CreateBoardAsEmoji() string {
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
	buttons := make([][]telebot.InlineButton, board.MaxCellSize)
	for r := range board.MaxCellSize {
		buttons[r] = make([]telebot.InlineButton, board.MaxCellSize)
		for c := range board.MaxCellSize {
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
				Data: string(gametype.XOGameType3X3) + "_play_" + strconv.Itoa(r) + strconv.Itoa(c),
			}
		}
	}
	return buttons
}

func (game *XOGame) CreatePlayersInlineButton(humanPlayers []*Player, CurrentPlayerTurn int) [][]telebot.InlineButton {
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
