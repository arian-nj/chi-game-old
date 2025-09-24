package socket

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/arian-nj/chibazi/backend/internals/utils"
	"github.com/coder/websocket"
	"golang.org/x/time/rate"
	"google.golang.org/protobuf/proto"
)

type Socket struct {
	Conn *websocket.Conn

	limiter *rate.Limiter

	Ctx       context.Context
	Cancel    context.CancelFunc
	EventChan chan []byte
}

func NewSocketClient(conn *websocket.Conn) *Socket {
	ctx, cancel := context.WithCancel(context.Background())
	return &Socket{
		Conn:      conn,
		limiter:   rate.NewLimiter(rate.Every(time.Millisecond*100), 10),
		Ctx:       ctx,
		Cancel:    cancel,
		EventChan: make(chan []byte, 16),
	}
}

func (sc *Socket) Listen(r *http.Request) {
	utils.RunBackgroundTask(func() {
		sc.listen(r)
	})
}

func (sc *Socket) listen(r *http.Request) {
	defer sc.Cancel()
	for {
		ctx := sc.Ctx

		err := sc.limiter.Wait(ctx)
		if err != nil {
			slog.Error("limiter burst exceded", "error", err)
		}

		_, messageByte, err := sc.Conn.Read(ctx)
		if err != nil {
			closeCode := websocket.CloseStatus(err)
			if closeCode == websocket.StatusNormalClosure {
				slog.Info("connection closed normally", "addr", r.RemoteAddr)
			} else {
				slog.Error("failed to read from socket", "addr", r.RemoteAddr, "err", err)
			}
			return
		}
		sc.EventChan <- messageByte
	}
}

func (sc *Socket) SendMessage(message proto.Message) error {
	out, err := proto.Marshal(message)
	if err != nil {
		return fmt.Errorf("can't marshal proto message %w", err)
	}
	if sc.Conn == nil {
		return errors.New("conn is nil")
	}
	return sc.Conn.Write(sc.Ctx, websocket.MessageBinary, out)

}
