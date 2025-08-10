package api

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/arian-nj/chibazi/database"
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
		OriginPatterns: CORS_PATTERNS,
	})
	if err != nil {
		slog.Error("error accepting new connection", "err", err)
		return
	}
	defer conn.CloseNow()

	socketClient := socket.NewSocketClient(conn)
	socketClient.Listen(r)

	var sessionPlayer *gamesessions.SessionPlayer

	for _, SPlayer := range gameSession.Players {
		if SPlayer.TgID == tgUser.TgID {
			sessionPlayer = SPlayer
			break
		}
	}
	sessionPlayer.Socket = socketClient

	for {
		select {
		case <-socketClient.Ctx.Done():
			slog.Info("socket context cancelled", "addr", r.RemoteAddr)
			return
		case newEvent := <-sessionPlayer.Socket.EventChan:
			gameSession.MsgChnl <- gamesessions.NewSessionEvent(sessionPlayer, newEvent)

		}
	}
}

type messageHistoryOut struct {
	Messages []messageOut `json:"messages"`
}

type messageOut struct {
	ID       int    `json:"id"`
	Text     string `json:"text"`
	SenderID int    `json:"sender_id"`
}

func (app *ApiApplication) getChatHistoryHandler(w http.ResponseWriter, r *http.Request) {
	tgUser, err := ContextGetAuthenticatedUser(app.Queries, r)
	if err != nil {
		app.ServerError(w, r, err)
		return
	}
	gameSession, ok := app.AllSessions.Get(strconv.Itoa(tgUser.TgID))
	if !ok {
		app.NotFound(w, r)
	}
	// NOTE: in case of dynamic Limit and offset cap limit by 50
	allMessages, err := app.Queries.GetSessionMessages(context.Background(), database.GetSessionMessagesParams{
		SessionID: gameSession.ID,
		Limit:     50,
	})
	if err != nil {
		app.ServerError(w, r, err)
		return
	}

	messageLen := len(allMessages)
	mhOut := messageHistoryOut{
		Messages: make([]messageOut, messageLen),
	}
	for index, message := range allMessages {
		mhOut.Messages[messageLen-1-index] = messageOut{
			ID:       message.ID,
			Text:     message.Message,
			SenderID: message.PlayerID,
		}
	}

	err = response.JSON(w, http.StatusOK, &mhOut)
	if err != nil {
		app.ServerError(w, r, err)
		return
	}
}
