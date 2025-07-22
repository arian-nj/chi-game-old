package humanplayer

import (
	"strconv"
	"time"

	"github.com/arian-nj/chibazi/internals/keybul"
)

type HumanPlayer struct {
	TgID int
	Name string

	MessageID int

	SpentTime     time.Duration
	TurnStartedAt time.Time
}

func NewHumanPlayer(name string, tgID int) *HumanPlayer {
	if len(name) > 20 {
		name = name[:20]
		name += "..."
	}
	return &HumanPlayer{
		Name: keybul.EscapeReserved(name),
		TgID: tgID,
	}
}

func (p *HumanPlayer) SetMessageSig(messageID int) *HumanPlayer {
	p.MessageID = messageID
	return p
}
func (p *HumanPlayer) MessageSig() (string, int64) {
	return strconv.Itoa(p.MessageID), int64(p.TgID)
}

func (p *HumanPlayer) Recipient() string {
	return strconv.FormatInt(int64(p.TgID), 10)
}
