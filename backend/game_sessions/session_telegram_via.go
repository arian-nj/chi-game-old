package gamesessions

import (
	"fmt"
	"log/slog"

	"github.com/arian-nj/chibazi/backend/internals/commander"
	"github.com/arian-nj/chibazi/backend/internals/keybul"
	"gopkg.in/telebot.v4"
)

type SessionTelegramViaListener struct {
	Bot          *telebot.Bot
	ViaMessageId string // Via Bots
}

func NewSessionTelegramViaListener(bot *telebot.Bot, viaMessageID string) *SessionTelegramViaListener {
	return &SessionTelegramViaListener{
		Bot:          bot,
		ViaMessageId: viaMessageID,
	}
}

func (tg *SessionTelegramViaListener) MessageSig() (string, int64) {
	return tg.ViaMessageId, 0
}

func (tg *SessionTelegramViaListener) Update(command commander.Command) {
	switch c := command.(type) {
	case *WaitForPlayerCommand:
		err := tg.SendWaitPanel(c)
		if err != nil {
			slog.Error("can not send wait panel", "error", err)
		}
	}
}

func (tg *SessionTelegramViaListener) SendWaitPanel(wait *WaitForPlayerCommand) error {
	session := wait.Session
	gameData := session.GameState.GetGameData()
	creator := wait.Creator
	inlineKeyboard := keybul.CreateInlineKeyboard(
		keybul.JoinGameInlineButtons,
	)
	text := gameData.StartText + "\n\n" + gameData.RulesText + "\n\n🕹 بازیکن " + fmt.Sprintf("%s", creator.Name) + " منتظر حریفه"

	return keybul.EditMessage(tg.Bot, tg, text, inlineKeyboard)
}
