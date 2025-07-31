package socket

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/coder/websocket"
	"golang.org/x/time/rate"
)

type Socket struct {
	Conn    *websocket.Conn
	limiter *rate.Limiter
	Ctx     context.Context
	Cancel  context.CancelFunc
}

func NewSocketClient(conn *websocket.Conn) *Socket {
	ctx, cancel := context.WithCancel(context.Background())
	return &Socket{
		Conn:    conn,
		limiter: rate.NewLimiter(rate.Every(time.Millisecond*100), 10),
		Ctx:     ctx,
		Cancel:  cancel,
	}
}

type EventMessage string

type EventType string

type Event struct {
	Type EventType
	Data *EventMessage
}

func NewEvent(Etype EventType, data *EventMessage) *Event {
	return &Event{
		Type: Etype,
		Data: data,
	}
}

func (sc *Socket) ListenToSocket() (*Event, error) {
	ctx := sc.Ctx

	if err := sc.limiter.Wait(ctx); err != nil {
		return nil, err
	}

	_, messageByte, err := sc.Conn.Read(ctx)
	if err != nil {
		return nil, err
	}

	newEvent := &Event{}
	err = json.Unmarshal(messageByte, newEvent)
	return newEvent, err
}

func (sc *Socket) SendEvent(newEvent *Event) error {
	data_btye, err := json.Marshal(newEvent)
	if err != nil {
		slog.Error("can't marshal event %w", "error", err)
		return err
	}

	err = sc.Conn.Write(sc.Ctx, websocket.MessageBinary, data_btye)
	if err != nil {
		slog.Error("can not write event %w", "error", err)
		return err
	}
	return nil
}

func (sc *Socket) Write(text string) error {
	return sc.Conn.Write(sc.Ctx, websocket.MessageText, []byte(text))
}
