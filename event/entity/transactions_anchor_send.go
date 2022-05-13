package event_entity

type TransactionsAnchorSend struct {
	Anchor int32  `json:"anchor"`
	State  string `json:"state"`
}

func NewTransactionsAnchorSend(anchor int32, state string) TransactionsAnchorSend {
	return TransactionsAnchorSend{
		Anchor: anchor,
		State:  state,
	}
}
