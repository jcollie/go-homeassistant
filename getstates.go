package homeassistant

// getStatesCommand .
type getStatesCommand struct {
	ID   uint64 `json:"id"`
	Type string `json:"type"`
}

// SetID .
func (m *getStatesCommand) SetID(id uint64) {
	m.ID = id
}

// GetStates .
func (ha *Connection) GetStates(handler ResultHandler) (uint64, error) {
	var message = getStatesCommand{
		Type: "get_states",
	}

	return ha.sendMessage(handler, &message)
}
