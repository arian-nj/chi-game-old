package api

import "github.com/arian-nj/chibazi/internals/socket"

const FinderEventType socket.EventType = "finder"

const (
	FMAdded   socket.EventMessage = "added"
	FMFound   socket.EventMessage = "found"
	FMTimeout socket.EventMessage = "timeout"
	FMCancel  socket.EventMessage = "cancel"
)
