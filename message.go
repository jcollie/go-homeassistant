package homeassistant

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Message message
type Message struct {
	ID      uint64          `json:"id"`
	Type    string          `json:"type"`
	Version string          `json:"ha_version"`
	Message string          `json:"message"`
	Success bool            `json:"success"`
	Result  json.RawMessage `json:"result"`
	Event   Event           `json:"event"`
	Error   Error           `json:"error"`
}

// Error .
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Context context
type Context struct {
	ID       string `json:"id"`
	ParentID string `json:"parent_id"`
	UserID   string `json:"user_id"`
}

// State state
type State struct {
	EntityID    string                     `json:"entity_id"`
	State       string                     `json:"state"`
	Attributes  map[string]json.RawMessage `json:"attributes"`
	LastChanged string                     `json:"last_changed"`
	LastUpdated string                     `json:"last_updated"`
	Context     Context                    `json:"context"`
}

//Data .
type Data struct {
	EntityID string `json:"entity_id"`
	OldState State  `json:"old_state"`
	NewState State  `json:"new_state"`
}

// Event event
type Event struct {
	EventType string  `json:"event_type"`
	Data      Data    `json:"data"`
	Origin    string  `json:"origin"`
	TimeFired string  `json:"time_fired"`
	Context   Context `json:"context"`
}

func (event *Event) String() string {
	var b strings.Builder

	fmt.Fprintf(&b, "event_type:                      %s\n", event.EventType)
	fmt.Fprintf(&b, "time_fired:                      %s\n", event.TimeFired)
	fmt.Fprintf(&b, "origin:                          %s\n", event.Origin)
	fmt.Fprintf(&b, "data.entity_id:                  %s\n", event.Data.EntityID)
	fmt.Fprintf(&b, "data.old_state.entity_id:        %s\n", event.Data.OldState.EntityID)
	fmt.Fprintf(&b, "data.old_state.state:            %s\n", event.Data.OldState.State)
	fmt.Fprintf(&b, "data.old_state.last_changed:     %s\n", event.Data.OldState.LastChanged)
	fmt.Fprintf(&b, "data.old_state.last_updated:     %s\n", event.Data.OldState.LastUpdated)
	fmt.Fprintf(&b, "data.new_state.state:            %s\n", event.Data.NewState.State)
	return b.String()
}
