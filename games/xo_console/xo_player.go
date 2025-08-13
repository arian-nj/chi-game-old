package xoconsole

import (
	"strconv"
	"time"

	"github.com/arian-nj/chibazi/internals/keybul"
)

type XoPlayer struct {
	TgID int
	Name string

	MessageID int

	SpentTime     time.Duration
	TurnStartedAt time.Time
}

func NewXoPlayer(name string, tgID int) *XoPlayer {
	if len(name) > 20 {
		name = name[:20]
		name += "..."
	}
	return &XoPlayer{
		Name: keybul.EscapeReserved(name),
		TgID: tgID,
	}
}

func (p *XoPlayer) SetMessageSig(messageID int) *XoPlayer {
	p.MessageID = messageID
	return p
}
func (p *XoPlayer) MessageSig() (string, int64) {
	return strconv.Itoa(p.MessageID), int64(p.TgID)
}

func (p *XoPlayer) Recipient() string {
	return strconv.FormatInt(int64(p.TgID), 10)
}
