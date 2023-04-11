package event_entity

type KeyUsageCounter struct {
	IssuerID      string `json:"issuer_id"`
	KeyProtection int    `json:"key_protection"`
}

func NewKeyUsageCounter(issuerID string, keyProtection int) KeyUsageCounter {
	return KeyUsageCounter{
		IssuerID:      issuerID,
		KeyProtection: keyProtection,
	}
}
