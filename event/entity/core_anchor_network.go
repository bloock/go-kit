package event_entity

type CoreAnchorNetwork struct {
	AnchorId  int    `json:"anchor_id"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	Test      bool   `json:"test"`
	TxHash    string `json:"tx_hash"`
	CreatedAt int64  `json:"created_at"`
}
