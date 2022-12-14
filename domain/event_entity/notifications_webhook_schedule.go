package event_entity

import "time"

type NotificationsWebhookSchedule struct {
	WType     string      `json:"type"`
	Payload   interface{} `json:"payload"`
	RequestId string      `json:"request_id"`
	Timestamp time.Time   `json:"timestamp"`
}

func NewNotificationsWebhookScheduleEventEntity(wType string, payload interface{}, rId string, timestamp time.Time) NotificationsWebhookSchedule {
	return NotificationsWebhookSchedule{
		WType:     wType,
		Payload:   payload,
		RequestId: rId,
		Timestamp: timestamp,
	}
}
