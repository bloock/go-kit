package event_entity

type TransactionsAnchorConfirm struct {
	AnchorID        int32  `json:"anchor_id"`
	TransactionHash string `json:"transaction_hash"`
	NetworkName     string `json:"network_name"`
}

func NewTransactionsAnchorConfirm(anchorID int32, txHash string, network string) TransactionsAnchorConfirm {
	return TransactionsAnchorConfirm{
		AnchorID:        anchorID,
		TransactionHash: txHash,
		NetworkName:     network,
	}
}
