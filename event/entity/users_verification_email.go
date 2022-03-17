package event_entity

type UsersVerificationEmail struct {
	Email       string `json:"email"`
	OriginURL   string `json:"origin_url"`
	VerifyToken string `json:"verify_token"`
}

func NewVerificationEmailEventEntity(email, url, verify string) UsersVerificationEmail {
	return UsersVerificationEmail{
		Email:       email,
		OriginURL:   url,
		VerifyToken: verify,
	}
}
