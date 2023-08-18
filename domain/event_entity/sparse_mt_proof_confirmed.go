package event_entity

type SparseMerkleTreeProofConfirmed struct {
	CredentialID string                `json:"credentialID"`
	Proof        SparseMerkleTreeProof `json:"proof"`
}

type SparseMerkleTreeProof struct {
	Type       string     `json:"type"`
	IssuerData IssuerData `json:"issuerData"`
	CoreClaim  string     `json:"coreClaim"`
	MTP        MTProof    `json:"mtp"`
}

type IssuerData struct {
	ID               string          `json:"id,omitempty"`
	State            IssuerDataState `json:"state,omitempty"`
	AuthCoreClaim    string          `json:"authCoreClaim,omitempty"`
	MTP              MTProof         `json:"mtp,omitempty"`
	CredentialStatus interface{}     `json:"credentialStatus,omitempty"`
}

type IssuerDataState struct {
	TxID               string `json:"txId,omitempty"`
	BlockTimestamp     int    `json:"blockTimestamp,omitempty"`
	BlockNumber        int    `json:"blockNumber,omitempty"`
	RootOfRoots        string `json:"rootOfRoots,omitempty"`
	ClaimsTreeRoot     string `json:"claimsTreeRoot,omitempty"`
	RevocationTreeRoot string `json:"revocationTreeRoot,omitempty"`
	Value              string `json:"value,omitempty"`
	Status             string `json:"status,omitempty"`
}

type MTProof struct {
	Existence  bool     `json:"existence"`
	Depth      uint     `json:"depth"`
	NotEmpties []byte   `json:"notempties"`
	Siblings   [][]byte `json:"siblings"`
	NodeAux    NodeAux  `json:"nodeAux"`
}

type NodeAux struct {
	Key   []byte `json:"key"`
	Value []byte `json:"value"`
}

func NewSparseMerkleTreeProofConfirmed(credentialID string, proof SparseMerkleTreeProof) SparseMerkleTreeProofConfirmed {
	return SparseMerkleTreeProofConfirmed{
		CredentialID: credentialID,
		Proof:        proof,
	}
}

func NewSparseMerkleTreeProof(_type string, coreClaim string, issuerData IssuerData, mtp MTProof) SparseMerkleTreeProof {
	return SparseMerkleTreeProof{
		Type:       _type,
		CoreClaim:  coreClaim,
		IssuerData: issuerData,
		MTP:        mtp,
	}
}

func NewIssuerData(id string, authCoreClaim string, credentialStatus interface{}, state IssuerDataState, mtp MTProof) IssuerData {
	return IssuerData{
		ID:               id,
		State:            state,
		AuthCoreClaim:    authCoreClaim,
		MTP:              mtp,
		CredentialStatus: credentialStatus,
	}
}

func NewIssuerDataState(txId string, blockTimestamp, blockNumber int, ror, ctr, rtr string, value, status string) IssuerDataState {
	return IssuerDataState{
		TxID:               txId,
		BlockTimestamp:     blockTimestamp,
		BlockNumber:        blockNumber,
		RootOfRoots:        ror,
		ClaimsTreeRoot:     ctr,
		RevocationTreeRoot: rtr,
		Value:              value,
		Status:             status,
	}
}

func NewNodeAux(key, value []byte) NodeAux {
	return NodeAux{
		Key:   key,
		Value: value,
	}
}

func NewMtProof(existence bool, depth uint, notEmpties []byte, siblings [][]byte, nodeAux NodeAux) MTProof {
	return MTProof{
		Existence:  existence,
		Depth:      depth,
		NotEmpties: notEmpties,
		Siblings:   siblings,
		NodeAux:    nodeAux,
	}
}
