package event_entity

type SparseMerkleTreeProofConfirmed struct {
	CredentialID string `json:"credential_id"`
	Proof        string `json:"proof"`
}

func NewSparseMerkleTreeProofConfirmed(credentialID string, proofEncoded string) SparseMerkleTreeProofConfirmed {
	return SparseMerkleTreeProofConfirmed{
		CredentialID: credentialID,
		Proof:        proofEncoded,
	}
}
