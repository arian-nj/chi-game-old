package xo

import (
	"fmt"
	"strconv"

	"github.com/arian-nj/chibazi/backend/internals/chronos"
	"github.com/arian-nj/chibazi/backend/internals/socket"
)

type XoPlayer struct {
	ID         int
	TelegramID int
	Name       string

	MessageID int

	Socket *socket.Socket

	// SpentTime         time.Duration
	// LastTurnStartedAt time.Time

	Timer *chronos.Timer
	Move  Cell
}

func NewXoPlayer(id int, name string, tgID int, socket *socket.Socket) *XoPlayer {
	if len(name) > 20 {
		name = name[:20]
		name += "..."
	}
	return &XoPlayer{
		// Name: keybul.EscapeReserved(name),
		ID:         id,
		Name:       fmt.Sprintf("`%s`", name),
		TelegramID: tgID,
		Timer:      chronos.NewTimer(MaxAllowedTime),
	}
}

func (p *XoPlayer) MessageSig() (string, int64) {
	return strconv.Itoa(p.MessageID), int64(p.TelegramID)
}

func (p *XoPlayer) Recipient() string {
	return strconv.FormatInt(int64(p.TelegramID), 10)
}
