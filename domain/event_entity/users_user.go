package event_entity

type UsersUser struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	VerifyToken string `json:"verify_token"`
	Business    string `json:"business"`
	Activated   bool   `json:"activated"`
	Country     string `json:"country"`
	VatID       string `json:"vat_id"`
	Deleted     bool   `json:"deleted"`
}

func NewUserEventEntity(id, name, surname, email, password, verifyToken, business string, activated, deleted bool, country string, vatID string) UsersUser {
	return UsersUser{
		ID:          id,
		Name:        name,
		Surname:     surname,
		Email:       email,
		Password:    password,
		VerifyToken: verifyToken,
		Business:    business,
		Country:     country,
		VatID:       vatID,
		Activated:   activated,
		Deleted:     deleted,
	}
}
