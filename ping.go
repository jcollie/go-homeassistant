package homeassistant

// pingCommand .
type pingCommand struct {
	ID   uint64 `json:"id"`
	Type string `json:"type"`
}

// SetID .
func (m *pingCommand) SetID(id uint64) {
	m.ID = id
}

// Ping .
func (ha *Connection) Ping(handler PongHandler) (uint64, error) {
	var message = pingCommand{
		Type: "ping",
	}

	return ha.sendMessage(handler, &message)
}
