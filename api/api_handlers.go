package api

import (
	"log/slog"
	"net/http"

	"github.com/arian-nj/chibazi/pkg/response"
	"github.com/coder/websocket"
)

func (app *ApiApplication) statusHandler(w http.ResponseWriter, r *http.Request) {
	err := response.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
	if err != nil {
		slog.Error(err.Error())
	}
}

func (app *ApiApplication) websocketUpgrader(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, nil) // this right error to w in case of error
	if err != nil {
		slog.Error("error accepting new connection", "err", err)
		return
	}
	defer conn.CloseNow()

	sokcetClient := NewSocketClient(conn)
	sokcetClient.Listen(r)
}
