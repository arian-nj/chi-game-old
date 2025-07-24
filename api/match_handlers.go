package api

import (
	"log/slog"
	"net/http"
	"time"

	gametype "github.com/arian-nj/chibazi/internals/game_type"
	matchmaking "github.com/arian-nj/chibazi/match_making"
	"github.com/coder/websocket"
)

func (app *ApiApplication) makeMatchMakingTicket(w http.ResponseWriter, r *http.Request) {
	tgUser, err := ContextGetAuthenticatedUser(app.Queries, r)
	if err != nil {
		app.ServerError(w, r, err)
		return
	}

	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		slog.Error("error accepting new connection", "err", err)
		return
	}
	defer conn.CloseNow()

	sokcetClient := NewSocketClient(conn)

	NewTicket := matchmaking.NewTicket("Player Name", tgUser.TgID, gametype.XOGameType3X3)
	app.MatchMaking.AddTicket(NewTicket)

	ticker := time.NewTicker(time.Second * 30)
	select {
	case <-NewTicket.MatchFound:
		sokcetClient.Write("found")
	case <-ticker.C:
		sokcetClient.Write("timeout")
	}
}
