package event_entity

type SubscriptionsSubscription struct {
	UserID    string   `json:"user_id"`
	Status    string   `json:"status"`
	ProductID []string `json:"product_id"`
}

func NewSubscriptionEventEntity(userID, status string, productID []string) SubscriptionsSubscription {
	return SubscriptionsSubscription{
		UserID:    userID,
		Status:    status,
		ProductID: productID,
	}
}
