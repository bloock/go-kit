package event_entity

type SubscriptionsTrialEnd struct {
	UserID         string `json:"user_id"`
	ExpirationDate string `json:"expiration_date"`
}
