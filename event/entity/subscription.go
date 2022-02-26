package event_entity

type Subscription struct {
	UserID string `json:"user_id"`
	Status string `json:"status"`
	PlanID string `json:"plan_id"`
}

func NewSubscriptionEventEntity(userID, status, planID string) Subscription {
	return Subscription{
		UserID: userID,
		Status: status,
		PlanID: planID,
	}
}
