package event_entity

type CredentialsPasswordResetEmail struct {
	Email      string `json:"email"`
	ResetToken string `json:"reset_token"`
}

func NewPasswordResetEmailEventEntity(email, resetToken string) CredentialsPasswordResetEmail {
	return CredentialsPasswordResetEmail{
		Email:      email,
		ResetToken: resetToken,
	}
}
