package events

import "fmt"

type EventType int32
type EventTypeTest int32

const (
	GetAnchor EventType = iota
	GetProof
	NewMessages
)

const(
	GetAnchorTest EventTypeTest = iota
	GetProofTest
	NewMessagesTest
)

const(
	CorePrefix     = "core"
	CoreTestPrefix = "core-test"
)

func (e EventType) Str() string {
	switch e {
	case GetAnchor:
		return fmt.Sprintf("%s.%s", CorePrefix, "get_anchor")
	case GetProof:
		return fmt.Sprintf("%s.%s", CorePrefix, "get_proof")
	case NewMessages:
		return fmt.Sprintf("%s.%s", CorePrefix, "new_messages")
	default:
		return ""
	}
}

func (e EventTypeTest) Str() string {
	switch e {
	case GetAnchorTest:
		return fmt.Sprintf("%s.%s", CoreTestPrefix, "get_anchor")
	case GetProofTest:
		return fmt.Sprintf("%s.%s", CoreTestPrefix, "get_proof")
	case NewMessagesTest:
		return fmt.Sprintf("%s.%s", CoreTestPrefix, "new_messages")
	default:
		return ""
	}
}
