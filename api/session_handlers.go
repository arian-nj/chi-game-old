package api

import (
	"log/slog"
	"net/http"
	"strconv"

	gamesessions "github.com/arian-nj/chibazi/game_sessions"
	"github.com/arian-nj/chibazi/internals/socket"
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
		return
	}

	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		slog.Error("error accepting new connection", "err", err)
		return
	}
	defer conn.CloseNow()

	socketClient := socket.NewSocketClient(conn)

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
		default:
			// newEvent, err := socketClient.Listen(r)
			// if err != nil {
			// 	if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
			// 		slog.Info("connection closed normally", "addr", r.RemoteAddr)
			// 	} else {
			// 		slog.Error("failed to read from socket", "addr", r.RemoteAddr, "err", err)
			// 	}
			// 	socketClient.Cancel()
			// 	slog.Error("error listening ", "error", err)
			// 	return
			// }
			//
			// newSessionEvent := gamesessions.NewSessionEvent(sessionPlayer, newEvent)
			// slog.Info("new event")
			// gameSession.MsgChnl <- newSessionEvent
		}
	}

}
