package event_entity

type UsersUser struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	VerifyToken string `json:"verify_token"`
	Country     string `json:"country"`
}

func NewUserEventEntity(id, name, surname, email, password, verifyToken string, country string) UsersUser {
	return UsersUser{
		ID:          id,
		Name:        name,
		Surname:     surname,
		Email:       email,
		Password:    password,
		VerifyToken: verifyToken,
		Country:     country,
	}
}
