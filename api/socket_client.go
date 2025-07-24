package api

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/coder/websocket"
	"golang.org/x/time/rate"
)

type SocketClient struct {
	Conn    *websocket.Conn
	limiter *rate.Limiter
	Ctx     context.Context
	cancel  context.CancelFunc
}

func NewSocketClient(conn *websocket.Conn) *SocketClient {
	ctx, cancel := context.WithCancel(context.Background())
	return &SocketClient{
		Conn:    conn,
		limiter: rate.NewLimiter(rate.Every(time.Millisecond*100), 10),
		Ctx:     ctx,
		cancel:  cancel,
	}
}
func (sc *SocketClient) Listen(r *http.Request) {
	for {
		err := sc.ListenToSocket()
		if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
			slog.Info("connection closed", "addr", r.RemoteAddr)
		}
		if err != nil {
			slog.Error("failed to echo", "addr", r.RemoteAddr, "err", err)
			sc.cancel()
			return
		}
	}

}
func (sc *SocketClient) ListenToSocket() error {
	ctx := context.Context(sc.Ctx) // FIXME:

	err := sc.limiter.Wait(ctx)
	if err != nil {
		return err
	}

	_, messageByte, err := sc.Conn.Read(ctx)
	if err != nil {
		return err
	}

	text := string(messageByte)
	slog.Info("got message", "text", text)

	return nil
}

func (sc *SocketClient) Write(text string) error {
	return sc.Conn.Write(sc.Ctx, websocket.MessageText, []byte(text))
}
