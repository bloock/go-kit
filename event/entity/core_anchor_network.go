package event_entity

type CoreAnchorNetwork struct {
	AnchorId  int    `json:"anchor_id"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	Test      bool   `json:"test"`
	TxHash    string `json:"tx_hash"`
	CreatedAt int64  `json:"created_at"`
}

func NewCoreAnchorNetworkEventEntity(anchorID int, name, status string, test bool, txHash string, createdAt int64) CoreAnchorNetwork {
	return CoreAnchorNetwork{
		AnchorId:  anchorID,
		Name:      name,
		Status:    status,
		Test:      test,
		TxHash:    txHash,
		CreatedAt: createdAt,
	}
}
