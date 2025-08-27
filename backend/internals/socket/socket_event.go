package socket

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
