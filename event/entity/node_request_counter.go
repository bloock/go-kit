package event_entity

type RequestsCounter struct {
	Counter int    `json:"counter"`
	UserID  string `json:"user_id"`
}

func NewRequestsCounterEventEntity(counter int, userID string) RequestsCounter {
	return RequestsCounter{
		Counter: counter,
		UserID:  userID,
	}
}
