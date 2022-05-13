package event_entity

type TransactionsAnchorSend struct {
	Anchor     int32  `json:"anchor"`
	State      string `json:"state"`
	Expiration int    `json:"expiration"`
}

func NewTransactionsAnchorSend(anchor int32, state string, expiration int) TransactionsAnchorSend {
	return TransactionsAnchorSend{
		Anchor:     anchor,
		State:      state,
		Expiration: expiration,
	}
}
