package homeassistant

import (
	"encoding/json"
	"time"
)

// Command .
type Command interface {
	SetID(id uint64)
}

// CloseHandler .
type CloseHandler interface {
	HandleClose(ha *Connection, id uint64)
}

// PongHandler defines an interface .
type PongHandler interface {
	HandlePong(ha *Connection, id uint64)
}

// ResultHandler .
type ResultHandler interface {
	CloseHandler
	HandleResult(ha *Connection, id uint64, success bool, result json.RawMessage)
}

// EventHandler .
type EventHandler interface {
	CloseHandler
	HandleEvent(ha *Connection, id uint64, origin string, timeFired time.Time, event *Event)
}

// ResultEventHandler .
type ResultEventHandler interface {
	ResultHandler
	EventHandler
}
