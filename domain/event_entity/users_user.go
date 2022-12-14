package event_entity

type UsersUser struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	VerifyToken string `json:"verify_token"`
	Activated   bool   `json:"activated"`
	Deleted     bool   `json:"deleted"`
}

func NewUserEventEntity(id, name, surname, email, password, verifyToken string, activated, deleted bool) UsersUser {
	return UsersUser{
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
