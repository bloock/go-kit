package event_entity

type NotificationsWebhookConfirmed struct {
	WebhookId          string `json:"webhook_id"`
	RequestId          string `json:"request_id"`
	WType              string `json:"type"`
	UserId             string `json:"user_id"`
	RequestBody        string `json:"request_body"`
	ResponseStatusCode int    `json:"response_status_code"`
	ResponseBody       string `json:"response_body"`
	Method             string `json:"method"`
	Url                string `json:"url"`
	CreatedAt          int64  `json:"created_at"`
}

func NewNotificationsWebhookConfirmedEventEntity(wID, rID, wType, uID, requestBody string, responseStatus int, responseBody, method, url string, createdAt int64) NotificationsWebhookConfirmed {
	return NotificationsWebhookConfirmed{
		WebhookId:          wID,
		RequestId:          rID,
		WType:              wType,
		UserId:             uID,
		RequestBody:        requestBody,
		ResponseStatusCode: responseStatus,
		ResponseBody:       responseBody,
		Method:             method,
		Url:                url,
		CreatedAt:          createdAt,
	}
}
