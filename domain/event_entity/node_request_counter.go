package event_entity

type RequestsCounter struct {
	Counter int    `json:"counter"`
	UserID  string `json:"user_id"`
	Network string `json:"network"`
}

func NewRequestsCounterEventEntity(counter int, userID string, network string) RequestsCounter {
	return RequestsCounter{
		Counter: counter,
		UserID:  userID,
		Network: network,
	}
}
