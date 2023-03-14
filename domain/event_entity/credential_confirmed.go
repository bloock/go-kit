package event_entity

type CredentialConfirmed struct {
	Payload interface{}
}

func NewCredentialConfirmedEventEntity(payload interface{}) CredentialConfirmed {
	return CredentialConfirmed{
		Payload: payload,
	}
}
