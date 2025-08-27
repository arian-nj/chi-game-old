package api

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"connectrpc.com/connect"
	"github.com/arian-nj/chibazi/backend/database"
	gamesessions "github.com/arian-nj/chibazi/backend/game_sessions"
	accountv1 "github.com/arian-nj/chibazi/backend/gen/account/v1"
	sessionv1 "github.com/arian-nj/chibazi/backend/gen/session/v1"
	"github.com/arian-nj/chibazi/backend/internals/socket"
	"github.com/arian-nj/chibazi/backend/pkg/response"
	"github.com/coder/websocket"
	"google.golang.org/protobuf/proto"
)

var (
	ErrorUnauthenticated = connect.NewError(connect.CodeUnauthenticated, errors.New("unknown user"))
)

func (app *ApiApplication) gameSessionWS(w http.ResponseWriter, r *http.Request) {
	tgUser, err := ContextGetAuthenticatedUser(app.Queries, r)
	if err != nil {
		app.InvalidAuthenticationToken(w, r)
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

func (app *ApiApplication) GetChatHistory(
	ctx context.Context,
	req *connect.Request[sessionv1.GetChatHistoryRequest],
) (*connect.Response[sessionv1.GetChatHistoryResponse], error) {
	person := app.AuthenticateHeader(ctx, req.Header())
	if person == nil {
		return nil, ErrorUnauthenticated
	}
	gameSession, ok := app.AllSessions.Get(strconv.Itoa(person.TgID))
	if !ok {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("no game found"))
	}
	// NOTE: in case of dynamic Limit and offset cap limit by 50
	allMessages, err := app.Queries.GetSessionMessages(context.Background(), database.GetSessionMessagesParams{
		SessionID: gameSession.ID,
		Limit:     50,
	})
	if err != nil {
		slog.Error("can't get message history", "error", err)
		return nil, connect.NewError(connect.CodeUnknown, errors.New("internal"))
	}

	messageLen := len(allMessages)
	// mhOut := chatHistoryOut{
	// 	Messages: make([]messageOut, messageLen),
	// }
	response := &sessionv1.GetChatHistoryResponse{
		Messages: make([]*sessionv1.ChatMessage, messageLen),
	}
	for index, message := range allMessages {
		response.Messages[messageLen-1-index] = &sessionv1.ChatMessage{
			Id:       int64(message.ID),
			Text:     message.Message,
			PlayerId: int64(message.PlayerID),
		}
	}

	return connect.NewResponse(response), nil
}

func (app *ApiApplication) GetSessionOpponent(
	ctx context.Context,
	req *connect.Request[sessionv1.GetSessionOpponentRequest],
) (*connect.Response[sessionv1.GetSessionOpponentResponse], error) {
	person := app.AuthenticateHeader(ctx, req.Header())
	if person == nil {
		return nil, ErrorUnauthenticated
	}

	gs, gsExist := app.AllSessions.Get(strconv.Itoa(person.TgID))
	if !gsExist {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("session not found"))
	}

	var opponents *gamesessions.SessionPlayer
	for _, p := range gs.Players {
		if p.ID != person.ID {
			opponents = p
			break
		}
	}

	return connect.NewResponse(&sessionv1.GetSessionOpponentResponse{
		Opponent: &accountv1.Account{
			Id:   int64(opponents.ID),
			Name: opponents.Name,
		},
	}), nil
}
