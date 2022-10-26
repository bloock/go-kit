package event_entity

type UsageRequest struct {
	UserID string `json:"user_id"`
	Value  int    `json:"value"`
}

func NewUsageRequest(userID string, value int) UsageRequest {
	return UsageRequest{
		UserID: userID,
		Value:  value,
	}
}
