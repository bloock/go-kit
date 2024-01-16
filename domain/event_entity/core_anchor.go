package event_entity

type CoreAnchor struct {
	EventType    string            `json:"event_type"`
	Id           int               `json:"id"`
	Root         string            `json:"root"`
	Finalized    bool              `json:"finalized"`
	Test         bool              `json:"test"`
	MessageCount int64             `json:"message_count"`
	CreatedAt    int64             `json:"created_at"`
	Network      CoreAnchorNetwork `json:"network"`
}

func NewCoreAnchorEventEntity(_type string, id int, root string, finalized, test bool, messageCount, createdAt int64, network CoreAnchorNetwork) CoreAnchor {
	return CoreAnchor{
		EventType:    _type,
		Id:           id,
		Root:         root,
		Finalized:    finalized,
		Test:         test,
		MessageCount: messageCount,
		CreatedAt:    createdAt,
		Network:      network,
	}
}
