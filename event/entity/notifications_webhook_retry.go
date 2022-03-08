package event_entity

type NotificationsWebhookRetry struct {
	Id         string `json:"id"`
	WebhookId  string `json:"webhook_id"`
	WType      string `json:"w_type"`
	Payload    string `json:"payload"`
	RequestId  string `json:"request_id"`
	Timestamp  string `json:"timestamp"`
	Expiration int    `json:"expiration"`
}

func NewWebhookRetryEventEntity(id, wId, wType, payload, rId, timestamp string, expiration int) NotificationsWebhookRetry {
	return NotificationsWebhookRetry{
		Id:         id,
		WebhookId:  wId,
		WType:      wType,
		Payload:    payload,
		RequestId:  rId,
		Timestamp:  timestamp,
		Expiration: expiration,
	}
}
