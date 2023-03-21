package event_entity

type KeyCreated struct {
	IssuerID      string `json:"issuer_id"`
	KeyID         string `json:"key_id"`
	KeyProtection int    `json:"key_protection"`
}

func NewKeyCreated(issuerID string, keyID string, keyProtection int) KeyCreated {
	return KeyCreated{
		IssuerID:      issuerID,
		KeyID:         keyID,
		KeyProtection: keyProtection,
	}
}
