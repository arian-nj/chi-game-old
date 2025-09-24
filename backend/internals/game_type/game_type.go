package gametype

type GameType string

var AllGameTypes = []GameType{XOGameType3X3, Conn4GameType}

const (
	XOGameType3X3 GameType = "xo3x3"
	XOGameType5X5 GameType = "xo5x5"
	Conn4GameType GameType = "conn4"
)

// DotBoxGameType    GameType = "dotbox"
// WebDotBoxGameType GameType = "webdotbox"
