package homeassistant

// getConfigCommand .
type getConfigCommand struct {
	ID   uint64 `json:"id"`
	Type string `json:"type"`
}

// SetID .
func (m *getConfigCommand) SetID(id uint64) {
	m.ID = id
}

// GetConfig .
func (ha *Connection) GetConfig(handler ResultHandler) (uint64, error) {
	var message = getConfigCommand{
		Type: "get_config",
	}

	return ha.sendMessage(handler, &message)
}
