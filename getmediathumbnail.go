package homeassistant

// getMediaPlayerThumbnailCommand .
type getMediaPlayerThumbnailCommand struct {
	ID       uint64 `json:"id"`
	Type     string `json:"type"`
	EntityID string `json:"entity_id"`
}

// SetID .
func (m *getMediaPlayerThumbnailCommand) SetID(id uint64) {
	m.ID = id
}

// GetMediaPlayerThumbnail .
func (ha *Connection) GetMediaPlayerThumbnail(entityID string, handler ResultHandler) (uint64, error) {
	var message = getMediaPlayerThumbnailCommand{
		Type:     "media_player_thumbnail",
		EntityID: entityID,
	}

	return ha.sendMessage(handler, &message)
}
