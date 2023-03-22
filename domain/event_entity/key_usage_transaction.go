package event_entity

type KeyUsageTransaction struct {
	Transacted int    `json:"transacted"`
	IssuerID   string `json:"issuer_id"`
}

func NewKeyUsageTransaction(transacted int, issuerID string) KeyUsageTransaction {
	return KeyUsageTransaction{
		Transacted: transacted,
		IssuerID:   issuerID,
	}
}
