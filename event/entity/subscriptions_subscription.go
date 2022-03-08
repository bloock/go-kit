package event_entity

type SubscriptionsSubscription struct {
	UserID string `json:"user_id"`
	Status string `json:"status"`
	PlanID string `json:"plan_id"`
}

func NewSubscriptionEventEntity(userID, status, planID string) SubscriptionsSubscription {
	return SubscriptionsSubscription{
		UserID: userID,
		Status: status,
		PlanID: planID,
	}
}
