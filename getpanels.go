package homeassistant

// getPanelsCommand .
type getPanelsCommand struct {
	ID   uint64 `json:"id"`
	Type string `json:"type"`
}

// SetID .
func (m *getPanelsCommand) SetID(id uint64) {
	m.ID = id
}

// GetPanels .
func (ha *Connection) GetPanels(handler ResultHandler) (uint64, error) {
	var message = getPanelsCommand{
		Type: "get_config",
	}

	return ha.sendMessage(handler, &message)
}
