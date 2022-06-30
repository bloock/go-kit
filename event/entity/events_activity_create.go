package event_entity

type EventsActivityCreate struct {
	Type         string `json:"type"`
	Status       int    `json:"status"`
	Path         string `json:"path"`
	RequestID    string `json:"request_id"`
	RequestBody  string `json:"request_body"`
	ResponseBody string `json:"response_body"`
	IP           string `json:"ip"`
	UserID       string `json:"user_id"`
}

func NewEventsActivityCreateEntity(typ string, st int, path, reqID, reqBody, respBody, ip, userID string) EventsActivityCreate {
	return EventsActivityCreate{
		Type:         typ,
		Status:       st,
		Path:         path,
		RequestID:    reqID,
		RequestBody:  reqBody,
		ResponseBody: respBody,
		IP:           ip,
		UserID:       userID,
	}
}
