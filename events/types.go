package events

type EventType int32

const (
	GetAnchor EventType = iota
	GetProof
	NewMessages
)

func (e EventType) Str() string {
	switch e {
	case GetAnchor:
		return "get_anchor"
	case GetProof:
		return "get_proof"
	case NewMessages:
		return "new_messages"
	default:
		return ""
	}
}
