package event_entity

type KeysWorker struct {
	Consumed int    `json:"consumed"`
	UserID   string `json:"user_id"`
}

func NewKeysWorker(consumed int, userID string) KeysWorker {
	return KeysWorker{
		Consumed: consumed,
		UserID:   userID,
	}
}
