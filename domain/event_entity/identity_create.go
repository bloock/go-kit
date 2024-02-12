package event_entity

type IdentityCreate struct {
	Did    string `json:"did"`
	KeyID  string `json:"key_id"`
	UserID string `json:"user_id"`
}
