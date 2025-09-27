package gamesessions

import (
	"fmt"
	"log/slog"

	"github.com/arian-nj/chibazi/backend/internals/commander"
	"github.com/arian-nj/chibazi/backend/internals/keybul"
	"gopkg.in/telebot.v4"
)

type SessionTelegramBotListener struct {
	Bot    *telebot.Bot
	UserID int
	TgID   int
}

func NewSessionTelegramBotListener(playerID int, TgID int, bot *telebot.Bot, viaMessageID string) *SessionTelegramBotListener {
	return &SessionTelegramBotListener{
		Bot:    bot,
		UserID: playerID,
		TgID:   TgID,
	}
}

func (tg *SessionTelegramBotListener) Update(command commander.Command) {
	switch c := command.(type) {
	case *MessageCommand:
		if c.Reciever.ID == tg.UserID {
			err := tg.SendChatMessageInBot(tg.TgID, c.Text, c.Sender.Name)
			if err != nil {
				slog.Error("Can not send chat message in session telegram", "error", err)
			}
		}

	case *GameEndedCommand:
		session := c.Session
		if session.Chat.IsOn == false {
			return
		}
		text := "چت قطع شد"
		_, err := session.Bot.Send(&telebot.User{ID: int64(tg.TgID)}, text, keybul.WelcomeReplyKeyboard)
		if err != nil {
			slog.Error("can't send chat ended message", "err", err)
		}

		text = fmt.Sprintf("چت تا %d ثانیه دیگه بسته میشه", int(ExpirationDur.Seconds()))
		_, err = c.Session.Bot.Send(&telebot.User{ID: int64(tg.TgID)}, text)
		if err != nil {
			slog.Error("can't send end game chat message", "err", err)
		}
	case *GameStartCommand:
		SendFoundOpponentMessage(c.Session.Players, tg.Bot)
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

func (tg *SessionTelegramBotListener) SendChatMessageInBot(toId int, text string, senderName string) error {
	_, err := tg.Bot.Send(&telebot.User{ID: int64(toId)},
		fmt.Sprintf("*_%s:_* %s", senderName, text), telebot.ModeMarkdownV2)
	return err
}
