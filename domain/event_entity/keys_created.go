package event_entity

type KeyCreated struct {
	IssuerID      string `json:"issuer_id"`
	KeyProtection int    `json:"key_protection"`
}

func NewKeyCreated(issuerID string, keyProtection int) KeyCreated {
	return KeyCreated{
		IssuerID:      issuerID,
		KeyProtection: keyProtection,
	}
}
