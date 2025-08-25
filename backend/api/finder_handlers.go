package api

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	finderv1 "github.com/arian-nj/chibazi/backend/gen/finder/v1"
	gametype "github.com/arian-nj/chibazi/backend/internals/game_type"
	"github.com/arian-nj/chibazi/backend/internals/socket"
	matchmaking "github.com/arian-nj/chibazi/backend/match_making"
	"github.com/coder/websocket"
	"google.golang.org/protobuf/proto"
)

func (app *ApiApplication) makeMatchMakingTicketWS(w http.ResponseWriter, r *http.Request) {
	tgUser, err := ContextGetAuthenticatedUser(app.Queries, r)
	if err != nil {
		app.InvalidAuthenticationToken(w, r)
		slog.Error("invalid auth", "err", err)
		return
	}

	hasTicket := app.MatchMaking.HasTicket(tgUser.TgID)
	if hasTicket {
		app.BadRequest(w, r, errors.New("user already have ticket can't make another one"))
		slog.Error("user already have ticket can't make another one")
		return
	}

	if app.AllSessions.IsSessionEmpty(tgUser.TgID) == false {
		app.BadRequest(w, r, errors.New("user already have active session can't make another one"))
		slog.Error("user already have ticket can't make another one")
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

	NewTicket := matchmaking.NewTicket("Player Name", tgUser.ID, tgUser.TgID, gametype.XOGameType3X3)
	app.MatchMaking.PushTicket(NewTicket)
	defer app.MatchMaking.RemovePlayerTicket(tgUser.TgID) // remove in case of error or canceling

	addEvent := finderv1.FinderEvent{
		Type: finderv1.FinderType_FINDER_TYPE_ADDED,
	}
	err = socketClient.SendMessage(&addEvent)
	if err != nil {
		slog.Error("", "error", err)
	}

	ticker := time.NewTicker(time.Second * 30)

	for {
		select {
		case <-NewTicket.MatchFoundChan:
			foundEvent := finderv1.FinderEvent{
				Type: finderv1.FinderType_FINDER_TYPE_FOUND,
			}
			socketClient.SendMessage(&foundEvent)
			return
		case <-ticker.C:
			timeoutEvent := finderv1.FinderEvent{
				Type: finderv1.FinderType_FINDER_TYPE_TIMEOUT,
			}
			socketClient.SendMessage(&timeoutEvent)
			return // removed by defer
		case newEventByte := <-socketClient.EventChan:
			newEvent := &finderv1.FinderEvent{}
			err := proto.Unmarshal(newEventByte, newEvent)
			if err != nil {
				slog.Error("can't unmarshal finder incoming event")
				continue
			}
			if newEvent.Type == finderv1.FinderType_FINDER_TYPE_CANCEL {
				slog.Info("Cancelling")
				return
			}
		case <-socketClient.Ctx.Done():
			return
		}
	}
}
