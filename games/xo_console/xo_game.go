package xoconsole

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/arian-nj/chibazi/database"
	gametype "github.com/arian-nj/chibazi/internals/game_type"
	humanplayer "github.com/arian-nj/chibazi/internals/human_player"
	keybul "github.com/arian-nj/chibazi/internals/keybul"
	"github.com/arian-nj/chibazi/internals/random"
	"github.com/arian-nj/chibazi/internals/xo_core"
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
const MaxPlayerTime = time.Minute * 2

type XOGame struct { // of GameInterface type
	GameType gametype.GameType

	XOBoard *xo_core.XoBoard

	ViaMessageId string // Via Bots
	LastEdit     time.Time

	Queries *database.Queries

	players []*humanplayer.HumanPlayer

	CancelGame context.CancelFunc
	Ctx        context.Context

	CurrentPlayerIndex int
}

func NewXOGame(gameType gametype.GameType, queries *database.Queries) *XOGame {
	maxBoardSize := 3
	winSize := 3
	if gameType == gametype.XOGameType5X5 {
		maxBoardSize = 5
		winSize = 4
	}
	randIndex := random.GenerateRandomNumber(2)

	ctx, cancel := context.WithCancel(context.Background())
	return &XOGame{

		CurrentPlayerIndex: randIndex,
		players:            []*humanplayer.HumanPlayer{},
		CancelGame:         cancel,
		Ctx:                ctx,

		GameType: gameType,
		XOBoard:  xo_core.NewTicBoard(maxBoardSize, winSize),
		Queries:  queries,
	}
}

func (cg *XOGame) MessageSig() (string, int64) {
	return cg.ViaMessageId, 0
}
func (cg *XOGame) Players() []*humanplayer.HumanPlayer {
	return cg.players
}
func (g *XOGame) AddPlayer(player *humanplayer.HumanPlayer) {
	g.players = append(g.players, player)
}

func (cg *XOGame) GetCurrentPlayer() *humanplayer.HumanPlayer {
	return cg.players[cg.CurrentPlayerIndex]
}

func (g *XOGame) NextPlayer() {
	if g.CurrentPlayerIndex == len(g.Players())-1 {
		g.CurrentPlayerIndex = 0
	} else {
		g.CurrentPlayerIndex += 1
	}
	g.GetCurrentPlayer().TurnStartedAt = time.Now()

}

func (cg *XOGame) GetContext() context.Context {
	return cg.Ctx
}

//
// func (g *XOGame) MonitorTimeout(bot telebot.API) {
// 	ticker := time.NewTicker(time.Second * 1)
// 	defer ticker.Stop()
// 	for {
// 		select {
// 		case <-ticker.C:
//
// 			now := time.Now()
// 			player := g.GetCurrentPlayer()
// 			player.SpentTime += now.Sub(player.TurnStartedAt)
// 			player.TurnStartedAt = time.Now()
//
// 			if player.SpentTime >= MaxPlayerTime {
// 				g.NextPlayer()
// 				err := g.TheEnd(bot, "\n برنده زمانی")
// 				if err != nil {
// 					slog.Error("error ending game with time out", "err", err)
// 				}
// 				return
// 			}
// 			if now.Sub(g.LastEdit) > time.Second*10 {
// 				err := g.EditDuringGameBoard(bot)
// 				if err != nil {
// 					slog.Error("can't edit message in time monitor", "err", err)
// 				}
//
// 			}
//
// 		case <-g.Ctx.Done():
// 			return
// 		}
// 	}
// }

func (g *XOGame) SendJoinPanelAddSender(c telebot.Context) error {
	sender := c.Sender()
	g.AddPlayer(humanplayer.NewHumanPlayer(sender.FirstName, int(sender.ID)))
	inlineKeyboard := keybul.CreateInlineKeyboard(
		keybul.JoinGameInlineButtons,
	)
	text := XOStartText + "\n\n" + g.RulesText() + "\n\n🕹 بازیکن " + fmt.Sprintf("%s", sender.FirstName) + " منتظر حریفه"
	return g.Edit(c.Bot(), g, text, inlineKeyboard)
}

func (g *XOGame) StartGame(bot telebot.API) error {
	err := g.Edit(bot, g, XOStartText+"\n\n"+g.RulesText(),
		keybul.CreateInlineKeyboard(
			keybul.CreateBotNameInlineButton(),
			CreateTicBoardInlineButton(g.XOBoard),
			g.CreatePlayersInlineButton(g.Players(), g.CurrentPlayerIndex),
		),
	)
	if err != nil {
		return fmt.Errorf("error when starting xo game %w", err)
	}
	for _, player := range g.Players() {
		now := time.Now()
		player.TurnStartedAt = now
	}
	// utils.RunBackgroundTask(func() {
	// 	g.MonitorTimeout(bot)
	// })
	return err
}

func (g *XOGame) CallbackHandler(c telebot.Context) error {
	callbackData := c.Callback().Data
	if callbackData == "join" {
		return g.XOJoinGameHandler(c)

	} else if after, hasPrefix := strings.CutPrefix(callbackData, "play_"); hasPrefix {
		// return g.XOPlayHandler(c, after)
	}
	return c.RespondAlert("no a valid callback")
}

func (g *XOGame) XOJoinGameHandler(c telebot.Context) error {
	sender := c.Callback().Sender
	if sender.ID == int64(g.Players()[0].TgID) {
		text := "خودت بازیو ساختی تو بازی هستی"
		return c.RespondText(text)
	}
	g.AddPlayer(humanplayer.NewHumanPlayer(sender.FirstName, int(sender.ID)))
	text := "اضافه شدی بازی شروع شد"
	err := c.RespondText(text)
	if err != nil {
		return err
	}
	return g.StartGame(c.Bot())
}

//
// func (g *XOGame) XOPlayHandler(c telebot.Context, callbackData string) error {
// 	sender := c.Sender()
// 	if len(callbackData) != 2 {
// 		return fmt.Errorf("invalid ttt_play data")
// 	}
//
// 	if g.GetCurrentPlayer().TgID != int(sender.ID) {
// 		return c.RespondText("نوبت تو نیست!")
// 	}
//
// 	xySlice := strings.Split(callbackData, "")
// 	rstr, cstr := xySlice[0], xySlice[1]
// 	rint, xerr := strconv.Atoi(rstr)
// 	cint, yerr := strconv.Atoi(cstr)
// 	if xerr != nil || yerr != nil {
// 		c.RespondAlert("یه مشکلی هست")
// 	}
//
// 	moveType := xo_core.Empty
// 	if g.CurrentPlayerIndex == 0 {
// 		moveType = xo_core.X
// 	} else {
// 		moveType = xo_core.O
// 	}
//
// 	isValid, errMsg := g.XOBoard.PlayMove(rint, cint, moveType)
// 	if !isValid {
// 		return c.RespondText(errMsg)
// 	}
//
// 	cellIndex := g.XOBoard.CellIndex(rint, cint)
// 	hasWon := g.XOBoard.HasWon(cellIndex)
// 	if hasWon {
// 		return g.TheEnd(c.Bot(), "")
// 	}
// 	if !g.XOBoard.IsAnyCellEmpty() {
// 		text := g.EndGameText() + "\nبازی مساوی شد"
//
// 		err := g.Edit(c.Bot(), g, text,
// 			keybul.CreateInlineKeyboard(
// 				keybul.CreateBotNameInlineButton(),
// 				keybul.EndGameInlineKeyboard(g.ViaMessageId != ""),
// 			),
// 		)
// 		return err
//
// 	}
// 	g.NextPlayer()
// 	return g.EditDuringGameBoard(c.Bot())
// }

// func (g *XOGame) EditDuringGameBoard(bot telebot.API) error {
// 	err := g.Edit(bot, g, XOStartText+"\n\n"+g.RulesText(),
// 		keybul.CreateInlineKeyboard(
// 			keybul.CreateBotNameInlineButton(),
// 			CreateTicBoardInlineButton(g.XOBoard),
// 			g.CreatePlayersInlineButton(g.Players(), g.CurrentPlayerIndex),
// 		),
// 	)
//
// 	return err
//
// }
//
// func (g *XOGame) TheEnd(bot telebot.API, additionalText string) error {
// 	g.CancelGame()
// 	text := g.EndGameText() + g.WinGameText() + additionalText
// 	err := g.Edit(bot, g, text,
// 		keybul.CreateInlineKeyboard(
// 			keybul.CreateBotNameInlineButton(),
// 			keybul.EndGameInlineKeyboard(g.ViaMessageId != ""),
// 		),
// 	)
// 	return err
// }
