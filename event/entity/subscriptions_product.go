package event_entity

type SubscriptionsProduct struct {
	ProductID   string                `json:"product_id"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Metadata    SubscriptionsMetadata `json:"metadata"`
}

type SubscriptionsMetadata struct {
	License     string `json:"license"`
	Optional    string `json:"optional"`
	Interval    string `json:"interval"`
	Plan        string `json:"plan"`
	ProductType string `json:"product_type"`
	Private     string `json:"private"`
}

func NewSubscriptionsMetadata(lic, op, in, pl, pt, prv string) SubscriptionsMetadata {
	return SubscriptionsMetadata{
		License:     lic,
		Optional:    op,
		Interval:    in,
		Plan:        pl,
		ProductType: pt,
		Private:     prv,
	}
}

func NewProductEventEntity(productID, name, description string, metadata SubscriptionsMetadata) SubscriptionsProduct {
	return SubscriptionsProduct{
		ProductID:   productID,
		Name:        name,
		Description: description,
		Metadata:    metadata,
	}
}
