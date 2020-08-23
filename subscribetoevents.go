package homeassistant

// SubscribeToEventsCommand .
type SubscribeToEventsCommand struct {
	ID        uint64 `json:"id"`
	Type      string `json:"type"`
	EventType string `json:"event_type"`
}

// SetID .
func (m *SubscribeToEventsCommand) SetID(id uint64) {
	m.ID = id
}

// SubscribeToEvents sends a command to Home Assistant to subscribe to events
func (ha *Connection) SubscribeToEvents(eventType string, handler ResultEventHandler) (uint64, error) {
	var message = SubscribeToEventsCommand{
		Type:      "subscribe_events",
		EventType: eventType,
	}

	return ha.sendMessage(handler, &message)
}
