package event_entity

type UsageRecord struct {
	UserID string `json:"user_id"`
	Value  int    `json:"value"`
}

func NewUsageRecord(userID string, value int) UsageRecord {
	return UsageRecord{
		UserID: userID,
		Value:  value,
	}
}
