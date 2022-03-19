package event_entity

type CoreAnchor struct {
	Id           int    `json:"id"`
	Root         string `json:"root"`
	Status       string `json:"status"`
	Test         bool   `json:"test"`
	MessageCount int64  `json:"message_count"`
	CreatedAt    int64  `json:"created_at"`
}

func NewCoreAnchorEventEntity(id int, root, status string, test bool, messageCount, createdAt int64) CoreAnchor {
	return CoreAnchor{
		Id:           id,
		Root:         root,
		Status:       status,
		Test:         test,
		MessageCount: messageCount,
		CreatedAt:    createdAt,
	}
}
