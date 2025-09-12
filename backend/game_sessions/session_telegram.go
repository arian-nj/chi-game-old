package gamesessions

import (
	"fmt"
	"log/slog"

	"github.com/arian-nj/chibazi/backend/internals/commander"
	"github.com/arian-nj/chibazi/backend/internals/keybul"
	"gopkg.in/telebot.v4"
)

type SessionTelegramListener struct {
	Bot          *telebot.Bot
	ViaMessageId string // Via Bots
	UserID       int
	TgID         int
}

func NewSessionTelegramListener(playerID int, TgID int, bot *telebot.Bot, viaMessageID string) *SessionTelegramListener {
	return &SessionTelegramListener{
		Bot:          bot,
		ViaMessageId: viaMessageID,
		UserID:       playerID,
		TgID:         TgID,
	}
}

func (tg *SessionTelegramListener) MessageSig() (string, int64) {
	return tg.ViaMessageId, 0
}

func (tg *SessionTelegramListener) Update(command commander.Command) {
	switch c := command.(type) {
	case *MessageCommand:
		if c.Reciever.ID == tg.UserID {
			err := tg.SendChatMessageInBot(tg.TgID, c.Text, c.Sender.Name)
			if err != nil {
				slog.Error("Can not send chat message in session telegram", "error", err)
			}
		}
	case *WaitForPlayerCommand:
		err := tg.SendWaitPanel(c)
		if err != nil {
			slog.Error("can not send wait panel", "error", err)
		}
	}
}

// Handler
func (session *GameSession) BotRequestSendMsg(bot telebot.API, senderID int, messageText string) error {
	if !session.Chat.IsOn {
		return nil
	}
	if len(messageText) > 256 {
		slog.Error("message is to long")
		return nil
	}

	var senderPlayer *SessionPlayer
	var recieverPlayer *SessionPlayer

	for _, p := range session.Players {
		if p.TgID == senderID {
			senderPlayer = p
		} else {
			recieverPlayer = p
		}
	}
	session.PushCommand(NewMessageCommand(session, messageText, senderPlayer, recieverPlayer))
	return nil
}

func (tg *SessionTelegramListener) SendChatMessageInBot(toId int, text string, senderName string) error {
	_, err := tg.Bot.Send(&telebot.User{ID: int64(toId)},
		fmt.Sprintf("*_%s:_* %s", senderName, text), telebot.ModeMarkdownV2)
	return err
}

func (tg *SessionTelegramListener) SendWaitPanel(wait *WaitForPlayerCommand) error {
	slog.Info("sending Wait Panel")
	session := wait.Session
	gameData := session.GameState.GetGameData()
	creator := wait.Creator
	inlineKeyboard := keybul.CreateInlineKeyboard(
		keybul.JoinGameInlineButtons,
	)
	text := gameData.StartText + "\n\n" + gameData.RulesText + "\n\n🕹 بازیکن " + fmt.Sprintf("%s", creator.Name) + " منتظر حریفه"

	return keybul.EditMessage(tg.Bot, tg, text, inlineKeyboard)
}
