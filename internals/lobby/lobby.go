package lobby

import (
	"github.com/arian-nj/chibazi/games/dotbox_console"
	xoconsole "github.com/arian-nj/chibazi/games/xo_console"
)

type Lobby struct {
	XOGames map[string]*xoconsole.XOGame
	DotBox  map[string]*dotbox_console.DotBoxGame
}
