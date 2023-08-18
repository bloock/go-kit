package event_entity

type IntegrityProofConfirmed struct {
	AnchorID int64  `json:"anchor_id"`
	Proof    string `json:"proof"`
}

func NewIntegrityProofConfirmed(anchorId int64, proof string) IntegrityProofConfirmed {
	return IntegrityProofConfirmed{
		AnchorID: anchorId,
		Proof:    proof,
	}
}
