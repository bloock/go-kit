package event_entity

type StorageWorker struct {
	Stored int    `json:"stored"`
	UserID string `json:"user_id"`
}

func NewStorageWorkerEventEntity(stored int, userID string) StorageWorker {
	return StorageWorker{
		Stored: stored,
		UserID: userID,
	}
}
