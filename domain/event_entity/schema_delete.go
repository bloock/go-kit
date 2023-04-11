package event_entity

type SchemaDelete struct {
	CID string `json:"CID"`
}

func NewSchemaDelete(CID string) SchemaDelete {
	return SchemaDelete{CID: CID}
}
