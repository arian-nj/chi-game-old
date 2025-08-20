package xoconsole

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/arian-nj/chibazi/internals/keybul"
	"gopkg.in/telebot.v4"
)

func (g *XOGame) CallBackRouter(c telebot.Context) error {
	callbackData := c.Callback().Data
	if callbackData == "join" {
		return g.XOJoinGameHandler(c)

	} else if after, hasPrefix := strings.CutPrefix(callbackData, "play_"); hasPrefix {
		return g.XOPlayHandler(c, after)
	}
	return c.RespondAlert("no a valid callback")
}

func (g *XOGame) XOJoinGameHandler(c telebot.Context) error {
	sender := c.Callback().Sender
	if sender.ID == int64(g.Players[0].TgID) {
		text := "خودت بازیو ساختی تو بازی هستی"
		return c.RespondText(text)
	}
	g.AddPlayer(sender.FirstName, int(sender.ID), nil)
	text := "اضافه شدی بازی شروع شد"
	err := c.RespondText(text)
	if err != nil {
		return err
	}
	return g.StartGame()
}

func (g *XOGame) StartGameBot() error {
	err := g.Edit(g.Bot, g, XOStartText+"\n\n"+g.RulesText(),
		keybul.CreateInlineKeyboard(
			keybul.CreateBotNameInlineButton(),
			CreateTicBoardInlineButton(g.Board),
			g.CreatePlayersInlineButton(g.Players, g.CurrentPlayerIndex),
		),
	)
	if err != nil {
		return err
	}

	return nil
}

func (g *XOGame) XOPlayHandler(c telebot.Context, callbackData string) error {
	sender := c.Sender()
	if len(callbackData) != 2 {
		return fmt.Errorf("invalid ttt_play data")
	}

	if g.GetCurrentPlayer().TgID != int(sender.ID) {
		return c.RespondText("نوبت تو نیست!")
	}

	xySlice := strings.Split(callbackData, "")
	rstr, cstr := xySlice[0], xySlice[1]
	rint, rerr := strconv.Atoi(rstr)
	cint, cerr := strconv.Atoi(cstr)
	if rerr != nil || cerr != nil {
		c.RespondAlert("یه مشکلی هست")
	}

	moveType := g.GetCurrentPlayer().Move

	isValid, errMsg := g.Board.ApplyMove(rint, cint, moveType)
	if !isValid {
		return c.RespondText(errMsg)
	}

	cellIndex := g.Board.CellIndex(rint, cint)
	g.SocketBrodcastNewMove(cellIndex, moveType)
	// hasWon := g.Board.HasWon(cellIndex)
	// if hasWon {
	// 	return g.TheEnd(c.Bot(), "")
	// }
	// if !g.Board.IsAnyCellEmpty() {
	// 	return g.TieGame(c.Bot())
	// }

	g.NextPlayer()
	return g.EditDuringGameBoard(c.Bot())
}

func (g *XOGame) EditDuringGameBoard(bot telebot.API) error {
	err := g.Edit(bot, g, XOStartText+"\n\n"+g.RulesText(),
		keybul.CreateInlineKeyboard(
			keybul.CreateBotNameInlineButton(),
			CreateTicBoardInlineButton(g.Board),
			g.CreatePlayersInlineButton(g.Players, g.CurrentPlayerIndex),
		),
	)
	return err
}
