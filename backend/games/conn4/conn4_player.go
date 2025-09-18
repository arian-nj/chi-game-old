package conn4

import (
	"fmt"
	"strconv"

	conn4_core "github.com/arian-nj/chibazi/backend/games/conn4/core"
	"github.com/arian-nj/chibazi/backend/internals/chronos"
	"github.com/arian-nj/chibazi/backend/internals/socket"
)

type Conn4Player struct {
	ID         int
	TelegramID int

	Name string

	MessageID int

	Socket *socket.Socket

	// SpentTime         time.Duration
	// LastTurnStartedAt time.Time

	Timer *chronos.Timer
	Move  conn4_core.Cell
}

func NewConn4Player(id int, name string, tgID int, socket *socket.Socket) *Conn4Player {
	if len(name) > 20 {
		name = name[:20]
		name += "..."
	}
	return &Conn4Player{
		// Name: keybul.EscapeReserved(name),
		ID:         id,
		Name:       fmt.Sprintf("`%s`", name),
		TelegramID: tgID,
		Timer:      chronos.NewTimer(MAX_ALLOWED_TIME),
	}
}

func (p *Conn4Player) MessageSig() (string, int64) {
	return strconv.Itoa(p.MessageID), int64(p.TelegramID)
}

func (p *Conn4Player) Recipient() string {
	return strconv.FormatInt(int64(p.TelegramID), 10)
}
