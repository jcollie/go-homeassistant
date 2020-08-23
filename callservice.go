package homeassistant

// callServiceCommand .
type callServiceCommand struct {
	ID          uint64      `json:"id"`
	Type        string      `json:"type"`
	Domain      string      `json:"domain"`
	Service     string      `json:"service"`
	ServiceData interface{} `json:"service_data"`
}

// SetID .
func (m *callServiceCommand) SetID(id uint64) {
	m.ID = id
}

// CallService calls a service
func (ha *Connection) CallService(domain string, service string, serviceData interface{}, handler ResultHandler) (uint64, error) {
	var message = callServiceCommand{
		Type:        "call_service",
		Domain:      domain,
		Service:     service,
		ServiceData: serviceData,
	}

	return ha.sendMessage(handler, &message)
}
