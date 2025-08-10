package api

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	gametype "github.com/arian-nj/chibazi/internals/game_type"
	"github.com/arian-nj/chibazi/internals/socket"
	matchmaking "github.com/arian-nj/chibazi/match_making"
	"github.com/coder/websocket"
)

func (app *ApiApplication) makeMatchMakingTicket(w http.ResponseWriter, r *http.Request) {
	tgUser, err := ContextGetAuthenticatedUser(app.Queries, r)
	if err != nil {
		app.ServerError(w, r, err)
		slog.Error("can't get user")
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
	app.MatchMaking.AddTicket(NewTicket)
	defer app.MatchMaking.RemovePlayerFromMatchMaking(tgUser.TgID)

	socketClient.SendNewEvent(FinderEventType, FMAdded)

	ticker := time.NewTicker(time.Second * 30)

	for {
		select {
		case <-NewTicket.MatchFoundChan:
			socketClient.SendNewEvent(FinderEventType, FMFound)
			return
		case <-ticker.C:
			socketClient.SendNewEvent(FinderEventType, FMTimeout)
		case newEvent := <-socketClient.EventChan:
			if newEvent.Type == FinderEventType && newEvent.Data == FMCancel {
				slog.Info("Cancelling")
				return
			}
		case <-socketClient.Ctx.Done():
			return
		}
	}
}
