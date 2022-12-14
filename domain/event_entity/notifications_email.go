package event_entity

type NotificationsEmail struct {
	From        string `json:"from"`
	To          string `json:"to"`
	Subject     string `json:"subject"`
	Body        string `json:"body"`
	AggregateID string `json:"aggregateID"`
	OccurredOn  string `json:"occurredOn"`
}
