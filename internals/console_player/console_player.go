package consoleplayer

import (
	"strconv"

	"github.com/arian-nj/chibazi/internals/keybul"
)

type ConsolePlayer struct {
	TgID int
	Name string

	MessageID int
}

func NewPlayer(name string, tgID int) *ConsolePlayer {
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
