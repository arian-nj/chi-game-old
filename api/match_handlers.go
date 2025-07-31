package api

import (
	"log/slog"
	"net/http"
	"time"

	gametype "github.com/arian-nj/chibazi/internals/game_type"
	"github.com/arian-nj/chibazi/internals/socket"
	matchmaking "github.com/arian-nj/chibazi/match_making"
	"github.com/coder/websocket"
)

func (app *ApiApplication) makeMatchMakingTicket(w http.ResponseWriter, r *http.Request) {
	slog.Info("hereee ")
	tgUser, err := ContextGetAuthenticatedUser(app.Queries, r)
	if err != nil {
		app.ServerError(w, r, err)
		slog.Error("can't get user")
		return
	}

	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		slog.Error("error accepting new connection", "err", err)
		return
	}
	defer conn.CloseNow()

	socketClient := socket.NewSocketClient(conn)

	NewTicket := matchmaking.NewTicket("Player Name", tgUser.TgID, gametype.XOGameType3X3)
	app.MatchMaking.AddTicket(NewTicket)

	ticker := time.NewTicker(time.Second * 30)

	select {
	case <-NewTicket.MatchFound:
		socketClient.Write("found")
	case <-ticker.C:
		socketClient.Write("timeout")
	}

	slog.Error("ended")
}
