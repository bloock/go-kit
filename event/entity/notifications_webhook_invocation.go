package event_entity

import "time"

type NotificationsWebhookInvocation struct {
	Id         string    `json:"id"`
	WebhookId  string    `json:"webhook_id"`
	WType      string    `json:"w_type"`
	Payload    interface{}    `json:"payload"`
	RequestId  string    `json:"request_id"`
	Timestamp  time.Time `json:"timestamp"`
	Expiration int       `json:"expiration"`
}

func NewWebhookInvocationEventEntity(id, wId, wType string, payload interface{}, rId string, timestamp time.Time, expiration int) NotificationsWebhookInvocation {
	return NotificationsWebhookInvocation{
		Id:         id,
		WebhookId:  wId,
		WType:      wType,
		Payload:    payload,
		RequestId:  rId,
		Timestamp:  timestamp,
		Expiration: expiration,
	}
}
