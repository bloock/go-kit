package event_entity

type IntegrityRecord struct {
	EventType    string            `json:"event_type"`
	Record       string            `json:"record"`
	ClientID     string            `json:"client_id"`
	AnchorID     int               `json:"anchor_id"`
	Test         bool              `json:"test"`
	CreatedAt    int64             `json:"created_at"`
	Networks     []CoreAnchorNetwork `json:"networks"`
}
