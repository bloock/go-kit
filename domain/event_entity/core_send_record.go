package event_entity

type SendRecord struct {
	Messages []string `json:"messages" binding:"required"`
}

func NewSendRecord(messages []string) *SendRecord {
	return &SendRecord{Messages: messages}
}
