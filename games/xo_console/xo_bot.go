package xoconsole

import (
	"github.com/arian-nj/chibazi/internals/keybul"
	"gopkg.in/telebot.v4"
)

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
