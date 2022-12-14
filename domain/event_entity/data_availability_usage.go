package event_entity

type DataAvailabilityUsage struct {
	UserID string `json:"user_id"`
	Value  int    `json:"value"`
}

func NewDataAvailabilityUsage(userID string, value int) DataAvailabilityUsage {
	return DataAvailabilityUsage{
		UserID: userID,
		Value:  value,
	}
}
