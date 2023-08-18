package event_entity

type SparseMerkleTreeProofConfirmed struct {
	CredentialID string      `json:"credential_id"`
	Proof        interface{} `json:"proof"`
}

func NewSparseMerkleTreeProofConfirmed(credentialID string, proof interface{}) SparseMerkleTreeProofConfirmed {
	return SparseMerkleTreeProofConfirmed{
		CredentialID: credentialID,
		Proof:        proof,
	}
}
