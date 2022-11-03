package event_entity

type TransferWorker struct {
	Transferred int    `json:"transferred"`
	UserID      string `json:"user_id"`
}

func NewTransferWorkerEventEntity(stored int, userID string) TransferWorker {
	return TransferWorker{
		Transferred: stored,
		UserID:      userID,
	}
}
