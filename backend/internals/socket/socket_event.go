package socket

import "encoding/json"

type SocketEventType string

type SocketEvent struct {
	Type SocketEventType `json:"type"`
	Data json.RawMessage `json:"data"`
}

func NewSocketEvent(Etype SocketEventType, data json.RawMessage) *SocketEvent {
	return &SocketEvent{
		Type: Etype,
		Data: data,
	}
}

func (e *SocketEvent) Decode(v any) error {
	return json.Unmarshal(e.Data, v)
}

type ActionType string

type GameAction struct {
	Action ActionType `json:"action"`
	AData  any        `json:"adata"`
}

func NewGameAction(action ActionType, data any) *GameAction {
	return &GameAction{
		Action: action,
		AData:  data,
	}
}
