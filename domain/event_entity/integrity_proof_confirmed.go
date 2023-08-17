package event_entity

type IntegrityProofConfirmed struct {
	Leaves []string    `json:"leaves"`
	Nodes  []string    `json:"nodes"`
	Depth  string      `json:"depth"`
	Bitmap string      `json:"bitmap"`
	Anchor ProofAnchor `json:"anchor"`
	Type   string      `json:"type"`
}

type ProofAnchor struct {
	AnchorID int64           `json:"anchor_id"`
	Networks []AnchorNetwork `json:"networks"`
	Root     string          `json:"root"`
	Status   string          `json:"status"`
}

type AnchorNetwork struct {
	Name   string `json:"name"`
	State  string `json:"state"`
	TxHash string `json:"tx_hash"`
}

func NewIntegrityProofConfirmed(leaves []string, nodes []string, depth string, bitmap string, _type string, anchor ProofAnchor) IntegrityProofConfirmed {
	return IntegrityProofConfirmed{
		Leaves: leaves,
		Nodes: nodes,
		Depth: depth,
		Bitmap: bitmap,
		Anchor: anchor,
		Type: _type,
	}
}

func NewProofAnchor(anchorId int64, root string, status string, networks []AnchorNetwork) ProofAnchor {
	return ProofAnchor {
		AnchorID: anchorId,
		Networks: networks,
		Root: root,
		Status: status,
	}
}

func NewAnchorNetwork(name string, state string, txHash string) AnchorNetwork {
	return AnchorNetwork {
		Name: name,
		State: state,
		TxHash: txHash,
	}
}
