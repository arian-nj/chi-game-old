package api

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	gamesessions "github.com/arian-nj/chibazi/game_sessions"
	"github.com/arian-nj/chibazi/internals/socket"
	"github.com/arian-nj/chibazi/pkg/response"
	"github.com/coder/websocket"
)

func (app *ApiApplication) gameSessionWS(w http.ResponseWriter, r *http.Request) {
	tgUser, err := ContextGetAuthenticatedUser(app.Queries, r)
	if err != nil {
		app.ServerError(w, r, err)
		return
	}

	gameSession, found := app.AllSessions.Get(strconv.Itoa(tgUser.TgID))
	if found == false {
		app.NotFound(w, r)
		slog.Info("not found")
		return
	}

	noconnct := r.URL.Query().Get("noconn")
	if noconnct != "" {
		response.JSON(w, http.StatusOK, nil)
		return
	}

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: COSS_PATTERNS,
	})
	if err != nil {
		slog.Error("error accepting new connection", "err", err)
		return
	}
	defer conn.CloseNow()

	socketClient := socket.NewSocketClient(conn)
	socketClient.Listen(r)

	var sessionPlayer *gamesessions.SessionPlayer

	for _, player := range gameSession.Players {
		if player.TgID == tgUser.TgID {
			sessionPlayer = player
			break
		}
	}
	sessionPlayer.Socket = socketClient

	for {
		select {
		case <-socketClient.Ctx.Done():
			slog.Info("socket context cancelled", "addr", r.RemoteAddr)
			return
		case <-sessionPlayer.Socket.EventChan:
			if newEvent.Type == FinderEventType {
				slog.Info("Cancelling")
				return
			}

		}
	}

}
