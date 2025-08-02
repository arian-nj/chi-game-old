package socket

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/arian-nj/chibazi/internals/utils"
	"github.com/coder/websocket"
	"golang.org/x/time/rate"
)

type Socket struct {
	Conn *websocket.Conn

	limiter *rate.Limiter

	Ctx       context.Context
	Cancel    context.CancelFunc
	EventChan chan *Event
}

func NewSocketClient(conn *websocket.Conn) *Socket {
	ctx, cancel := context.WithCancel(context.Background())
	return &Socket{
		Conn:      conn,
		limiter:   rate.NewLimiter(rate.Every(time.Millisecond*100), 10),
		Ctx:       ctx,
		Cancel:    cancel,
		EventChan: make(chan *Event, 16),
	}
}

type EventMessage string

type EventType string

type Event struct {
	Type EventType    `json:"type"`
	Data EventMessage `json:"data"`
}

func NewEvent(Etype EventType, data EventMessage) *Event {
	return &Event{
		Type: Etype,
		Data: data,
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
		newEvent := &Event{}
		err = json.Unmarshal(messageByte, newEvent)
		if err != nil {
			slog.Error("can't marshal websocket event", "error", err)
		} else {
			sc.EventChan <- newEvent
		}
	}
}

func (sc *Socket) SendNewEvent(Etype EventType, data EventMessage) error {
	return sc.SendEvent(NewEvent(Etype, data))
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
