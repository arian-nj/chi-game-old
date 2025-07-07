package dotbox_console

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/arian-nj/chibazi/database"
	commonapp "github.com/arian-nj/chibazi/internals/common_app"
	dotbox "github.com/arian-nj/chibazi/internals/dot_box"
	gametype "github.com/arian-nj/chibazi/internals/game_type"
	keybul "github.com/arian-nj/chibazi/internals/keybul"
	"github.com/arian-nj/chibazi/internals/random"
	"gopkg.in/telebot.v4"
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

type DotBoxGame struct { // of GameInterface type
	DotBoxBoard        *dotbox.DotBoxBoard
	CurrentPlayerIndex int
	Players            []*Player
	MessageId          string
	PlayerOnScore      int
	PlayerTwoScore     int

	CreatedAt time.Time
	common    *commonapp.CommonApp
}

func NewDotBoxGame(common *commonapp.CommonApp) *DotBoxGame {
	return &DotBoxGame{
		DotBoxBoard: dotbox.NewDotBoard(),
		Players:     []*Player{},
		CreatedAt:   time.Now(),
		common:      common,
	}
}

const (
	DotBoxStartText = `❌ *نقطه بازی* ⭕️`
)

func (g *DotBoxGame) SendJoinPanel(c telebot.Context) error {
	sender := c.Sender()
	g.addPlayer(newPlayer(sender.FirstName, int(sender.ID)))
	inlineKeyboard := keybul.CreateInlineKeyboard(
		keybul.JoinGameKeyboard(gametype.DotBoxGameType),
	)
	text := DotBoxStartText + "\n\n" + g.RulesText() + "\n\n🕹 بازیکن " + fmt.Sprintf("%s", sender.FirstName) + " منتظر حریفه"
	return keybul.EditGameMessage(c.Bot(), g, text, inlineKeyboard)
}

func (g *DotBoxGame) NextPlayer() {
	if g.CurrentPlayerIndex == len(g.Players)-1 {
		g.CurrentPlayerIndex = 0
		return
	}
	g.CurrentPlayerIndex += 1
}

func (g *DotBoxGame) MessageSig() (string, int64) {
	return g.MessageId, 0
}

func (g *DotBoxGame) addPlayer(player *Player) {
	g.Players = append(g.Players, player)
}
func (game *DotBoxGame) StartGame(c telebot.Context) error {
	randIndex := random.GenerateRandomNumber(2)
	game.CurrentPlayerIndex = randIndex

	err := keybul.EditGameMessage(c.Bot(), game, DotBoxStartText+"\n\n"+game.RulesText(),
		keybul.CreateInlineKeyboard(
			keybul.CreateBotNameInlineButton(),
			CreateDotBoxInlineButton(game.DotBoxBoard),
			game.CreatePlayersInlineButton(game.Players, game.CurrentPlayerIndex),
		),
	)
	if err != nil {
		return err
	}
	_, err = game.common.Queries.CreateHub(context.Background(), database.CreateHubParams{
		GameType: string(gametype.DotBoxGameType),
		TgID:     game.Players[0].TgID,
	})
	return err
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
				Data: string(gametype.DotBoxGameType) + "_play_" + strconv.Itoa(r) + strconv.Itoa(c),
			}
			row = append(row, btn)
		}
		buttons = append(buttons, row)
	}
	return buttons
}

func (game *DotBoxGame) CallbackHandlers(c telebot.Context, callbackData string) error {
	if strings.HasPrefix(callbackData, "play_") {
		return game.PlayHandler(c)
	} else if strings.HasPrefix(callbackData, "join") {
		return game.JoinGameHandler(c, callbackData)
	}
	return nil
}
func (game *DotBoxGame) JoinGameHandler(c telebot.Context, callbackData string) error {
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
func (game *DotBoxGame) PlayHandler(c telebot.Context) error {
	callback := c.Callback()
	callbackData := strings.TrimPrefix(callback.Data, string(gametype.DotBoxGameType)+"_play_")

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

		return keybul.EditGameMessage(c.Bot(), game, game.EndGameText()+text,
			keybul.CreateInlineKeyboard(
				keybul.CreateBotNameInlineButton(),
				keybul.EndgameInlineKeyboard,
			),
		)

	}
	return keybul.EditGameMessage(c.Bot(), game, DotBoxStartText+"\n\n"+game.RulesText()+"\n"+game.CreateBoardAsEmoji(), keybul.CreateInlineKeyboard(
		keybul.CreateBotNameInlineButton(),
		CreateDotBoxInlineButton(game.DotBoxBoard),
		game.CreatePlayersInlineButton(game.Players, game.CurrentPlayerIndex),
	))

}

func (game *DotBoxGame) EndGameText() string {
	return DotBoxStartText + "\nبازیکن ها:\n" + game.Players[0].Name + " " + strconv.Itoa(game.PlayerOnScore) + "\n" + game.Players[1].Name + " " + strconv.Itoa(game.PlayerTwoScore)
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

func (game *DotBoxGame) CreatePlayersInlineButton(humanPlayers []*Player, CurrentPlayerTurn int) [][]telebot.InlineButton {
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
			// URL:  fmt.Sprintf("tg://user?id=%d", hplayer.TgID),
			Data: "_",
		})
		buttons = append(buttons, row)
	}

	return buttons
}

func (game *DotBoxGame) RulesText() string {
	text := ""
	text += "هر دکمه گوشه یک مربعه\n"
	text += "هر مربعی که کامل کنی یک امتیاز میگیری\n"
	text += "اگه امتیاز بگیری نوبتت نمیگذره"
	return text
}
