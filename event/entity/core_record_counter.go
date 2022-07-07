package event_entity

type RecordsCounter struct {
	Counter int    `json:"counter"`
	UserID  string `json:"user_id"`
}

func NewRecordsCounterEventEntity(counter int, userID string) RecordsCounter {
	return RecordsCounter{
		Counter: counter,
		UserID:  userID,
	}
}
