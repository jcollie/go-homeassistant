package homeassistant

// getServicesCommand .
type getServicesCommand struct {
	ID   uint64 `json:"id"`
	Type string `json:"type"`
}

// SetID .
func (m *getServicesCommand) SetID(id uint64) {
	m.ID = id
}

// GetServices .
func (ha *Connection) GetServices(handler ResultHandler) (uint64, error) {
	var message = getServicesCommand{
		Type: "get_config",
	}

	return ha.sendMessage(handler, &message)
}
