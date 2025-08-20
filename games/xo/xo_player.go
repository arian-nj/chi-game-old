package xo

import (
	"fmt"
	"strconv"
	"time"

	"github.com/arian-nj/chibazi/internals/socket"
)

type XoPlayer struct {
	TgID int
	Name string

	MessageID int

	Socket *socket.Socket

	SpentTime     time.Duration
	TurnStartedAt time.Time

	Move Cell
}

func NewXoPlayer(name string, tgID int, socket *socket.Socket) *XoPlayer {
	if len(name) > 20 {
		name = name[:20]
		name += "..."
	}
	return &XoPlayer{
		// Name: keybul.EscapeReserved(name),
		Name: fmt.Sprintf("`%s`", name),
		TgID: tgID,
	}
}

func (p *XoPlayer) MessageSig() (string, int64) {
	return strconv.Itoa(p.MessageID), int64(p.TgID)
}

func (p *XoPlayer) Recipient() string {
	return strconv.FormatInt(int64(p.TgID), 10)
}
