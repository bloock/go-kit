package event_entity

type KeysUsage struct {
	UserID string `json:"user_id"`
	Value  int    `json:"value"`
}

func NewKeysUsage(userID string, value int) KeysUsage {
	return KeysUsage{
		UserID: userID,
		Value:  value,
	}
}
