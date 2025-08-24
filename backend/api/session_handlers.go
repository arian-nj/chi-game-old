package api

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/arian-nj/chibazi/backend/database"
	gamesessions "github.com/arian-nj/chibazi/backend/game_sessions"
	sessionv1 "github.com/arian-nj/chibazi/backend/gen/session/v1"
	"github.com/arian-nj/chibazi/backend/internals/socket"
	"github.com/arian-nj/chibazi/backend/pkg/response"
	"github.com/coder/websocket"
	"google.golang.org/protobuf/proto"
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
	if gameSession.GameState != nil {
		gameSession.GameState.SetPlayerSocket(tgUser.TgID, socketClient)
	}

	for {
		select {
		case <-socketClient.Ctx.Done():
			slog.Info("socket context cancelled", "addr", r.RemoteAddr)
			return
		case newMsgBytes := <-sessionPlayer.Socket.EventChan:
			newSessionMsg := &sessionv1.SessionMessage{}
			err := proto.Unmarshal(newMsgBytes, newSessionMsg)
			if err != nil {
				slog.Error("can't unmarshal session msg", "error", err)
				continue
			}
			gameSession.MsgChnl <- gamesessions.NewSessionEvent(sessionPlayer, newSessionMsg)
		}
	}
}

type chatHistoryOut struct {
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
		return
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
	mhOut := chatHistoryOut{
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
