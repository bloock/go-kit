package event_entity

type EventsActivityCreate struct {
	Type         string `json:"type"`
	Status       int    `json:"status"`
	Path         string `json:"path"`
	RequestID    string `json:"request_id"`
	RequestBody  string `json:"request_body"`
	ResponseBody string `json:"response_body"`
	IpAddress    string `json:"ip_address"`
	UserID       string `json:"user_id"`
	Method       string `json:"method"`
}

func NewEventsActivityCreateEntity(typ string, st int, path, reqID, reqBody, respBody, ip, userID, method string) EventsActivityCreate {
	return EventsActivityCreate{
		Type:         typ,
		Status:       st,
		Path:         path,
		RequestID:    reqID,
		RequestBody:  reqBody,
		ResponseBody: respBody,
		IpAddress:    ip,
		UserID:       userID,
		Method:       method,
	}
}
