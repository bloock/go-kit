package event_entity

type KeysTransactions struct {
	Transacted int    `json:"transacted"`
	IssuerID   string `json:"issuer_id"`
}

func NewKeysTransactions(transacted int, issuerID string) KeysTransactions {
	return KeysTransactions{
		Transacted: transacted,
		IssuerID:   issuerID,
	}
}
