package event_entity

type User struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	VerifyToken string `json:"verify_token"`
	Activated   bool   `json:"activated"`
	Deleted     bool   `json:"deleted"`
}

func NewUserEventEntity(id, name, surname, email, password, verifyToken string, activated, deleted bool) User {
	return User{
		ID:          id,
		Name:        name,
		Surname:     surname,
		Email:       email,
		Password:    password,
		VerifyToken: verifyToken,
		Activated:   activated,
		Deleted:     deleted,
	}
}
