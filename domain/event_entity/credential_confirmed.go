package event_entity

type CredentialConfirmed struct {
	ThreadID string                 `json:"thid"`
	Body     ClaimOfferBodyResponse `json:"body"`
	From     string                 `json:"from"`
	To       string                 `json:"to"`
}

type ClaimOfferBodyResponse struct {
	URL         string            `json:"url"`
	Credentials []CredentialOffer `json:"credentials"`
}

type CredentialOffer struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

func NewCredentialConfirmedEventEntity(threadID, from, to string, url string, id string, description string) CredentialConfirmed {
	credOffer := CredentialOffer{
		ID: id,
		Description: description,
	}

	body := ClaimOfferBodyResponse{
		URL: url,
		Credentials: []CredentialOffer{credOffer},
	}

	return CredentialConfirmed{
		ThreadID: threadID,
		Body: body,
		From: from,
		To: to,
	}
}
