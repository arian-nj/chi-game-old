package consoleplayer

import (
	"strconv"
	"time"

	"github.com/arian-nj/chibazi/internals/keybul"
)

type ConsolePlayer struct {
	TgID int
	Name string

	MessageID int

	SpentTime     time.Duration
	TurnStartedAt time.Time
}

func NewConsolePlayer(name string, tgID int) *ConsolePlayer {
	if len(name) > 20 {
		name = name[:20]
		name += "..."
	}
	return &ConsolePlayer{
		Name: keybul.EscapeReserved(name),
		TgID: tgID,
	}
}

func (p *ConsolePlayer) SetMessageSig(messageID int) *ConsolePlayer {
	p.MessageID = messageID
	return p
}
func (p *ConsolePlayer) MessageSig() (string, int64) {
	return strconv.Itoa(p.MessageID), int64(p.TgID)
}

func (p *ConsolePlayer) Recipient() string {
	return strconv.FormatInt(int64(p.TgID), 10)
}
