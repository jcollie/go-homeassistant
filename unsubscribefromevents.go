package homeassistant

// UnsubscribeFromEventsCommand .
type UnsubscribeFromEventsCommand struct {
	ID           uint64 `json:"id"`
	Type         string `json:"type"`
	Subscription uint64 `json:"subscription"`
}

// SetID .
func (m *UnsubscribeFromEventsCommand) SetID(id uint64) {
	m.ID = id
}

// UnsubscribeFromEvents sends a command to Home Assistant to unsubscribe to events
func (ha *Connection) UnsubscribeFromEvents(subscription uint64, handler EventHandler) (uint64, error) {
	var message = UnsubscribeFromEventsCommand{
		Type:         "unsubscribe_events",
		Subscription: subscription,
	}

	return ha.sendMessage(handler, &message)
}
