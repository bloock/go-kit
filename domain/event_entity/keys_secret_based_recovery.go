package event_entity

type KeysSecretBasedRecovery struct {
	Email string `json:"email"`
	Code  int    `json:"code"`
}
