package api

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/arian-nj/chibazi/pkg/response"
	"github.com/coder/websocket"
	"golang.org/x/time/rate"
)

func (app *Application) statusHandler(w http.ResponseWriter, r *http.Request) {
	err := response.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
	if err != nil {
		slog.Error(err.Error())
	}
}

func (app *Application) websocketUpgrader(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		Subprotocols: []string{"echo"},
	})
	if err != nil {
		slog.Error("error accepting new connection", "err", err)
		return
	}
	defer conn.CloseNow()

	if conn.Subprotocol() != "echo" {
		conn.Close(websocket.StatusPolicyViolation, "client must speak the echo subprotocol")
		return
	}

	limitter := rate.NewLimiter(rate.Every(time.Millisecond*100), 10)
	for {
		err := echo(conn, limitter)
		if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
			return
		}
		if err != nil {
			slog.Error("failed to echo", "addr", r.RemoteAddr, "err", err)
			return
		}
	}
}

func echo(c *websocket.Conn, l *rate.Limiter) error {
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	// defer cancel()

	ctx := context.Background()

	err := l.Wait(ctx)
	if err != nil {
		return err
	}

	typ, r, err := c.Reader(ctx)
	if err != nil {
		return err
	}

	w, err := c.Writer(ctx, typ)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, r)
	if err != nil {
		return fmt.Errorf("failed to copy: %w", err)
	}
	w.Close()

	return nil
}
