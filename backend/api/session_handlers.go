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
	"github.com/coder/websocket"
	"google.golang.org/protobuf/proto"
)

var (
	ErrorUnauthenticated = connect.NewError(connect.CodeUnauthenticated, errors.New("unknown user"))
)

func sendSessionSocketError(socketClient *socket.Socket, errType sessionv1.SessionErrorType) {
	addEvent := sessionv1.SessionMessage{
		Content: &sessionv1.SessionMessage_Error{
			Error: errType,
		},
	}
	err := socketClient.SendMessage(&addEvent)
	if err != nil {
		slog.Error("", "error", err)
	}
}

func (app *ApiApplication) sessionWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: CORS_PATTERNS,
	})
	if err != nil {
		slog.Error("error accepting new connection", "err", err)
		return
	}
	defer conn.Close(websocket.StatusNormalClosure, "websocket ended")

	socketClient := socket.NewSocketClient(conn)

	personRow, err := ContextGetAuthenticatedUser(app.Queries, r)
	if err != nil {
		sendSessionSocketError(socketClient, sessionv1.SessionErrorType_SESSION_ERROR_TYPE_AUTH)
		return
	}

	gameSession, found := app.AllSessions.Get(strconv.Itoa(personRow.ID))
	if found == false {
		sendSessionSocketError(socketClient, sessionv1.SessionErrorType_SESSION_ERROR_TYPE_NOSESSION)
		return
	}

	socketClient.Listen(r)

	var sessionPlayer *gamesessions.SessionPlayer

	for _, SPlayer := range gameSession.Players {
		if SPlayer.ID == personRow.ID {
			sessionPlayer = SPlayer
			break
		}
	}
	if sessionPlayer == nil {
		slog.Error("no player found in session", "id", personRow.ID)
		return
	}
	sessionPlayer.Socket = socketClient
	gamesessions.SendGametypeOverSocket(gameSession, sessionPlayer)

	socketSubber := gamesessions.NewSessionSocketListener(sessionPlayer)
	gameSession.Subscribe(socketSubber)
	defer gameSession.Unsubscribe(socketSubber)

	if gameSession.GameState != nil {
		cancel := gameSession.GameState.SubToSocket(personRow.ID, socketClient)
		if cancel == nil {
			sendSessionSocketError(socketClient, sessionv1.SessionErrorType_SESSION_ERROR_TYPE_UNSPECIFIED)
			return
		}
		defer cancel()
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
	gameSession, ok := app.AllSessions.Get(strconv.Itoa(person.ID))
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

	gs, gsExist := app.AllSessions.Get(strconv.Itoa(person.ID))
	if !gsExist {
		slog.Error("session not found", "person", person.ID)
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

func (app *ApiApplication) HasSession(
	ctx context.Context,
	req *connect.Request[sessionv1.HasSessionRequest],
) (*connect.Response[sessionv1.HasSessionResponse], error) {

	person := app.AuthenticateHeader(ctx, req.Header())
	if person == nil {
		return nil, ErrorUnauthenticated
	}
	_, hasSession := app.AllSessions.Get(strconv.Itoa(person.ID))

	return connect.NewResponse(&sessionv1.HasSessionResponse{
		HasSession: hasSession,
	}), nil
}
