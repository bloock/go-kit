package network

type NetworkName int32

const (
	EthMainnet NetworkName = iota
	GnosisChain
	PolygonChain
	EthBloockchain
	EthRinkeby
	EthRopsten
	EthGoerli
	EthSepolia
	NetworkDefault
)

func NewNetworkName(state string) NetworkName {
	switch state {
	case "bloock_chain":
		return EthBloockchain
	case "ethereum_mainnet":
		return EthMainnet
	case "ethereum_rinkeby":
		return EthRinkeby
	case "ethereum_ropsten":
		return EthRopsten
	case "gnosis_chain":
		return GnosisChain
	case "polygon_chain":
		return PolygonChain
	case "ethereum_goerli":
		return EthGoerli
	case "ethereum_sepolia":
		return EthSepolia
	default:
		return EthRinkeby
	}
}

func (n NetworkName) Str() string {
	switch n {
	case EthBloockchain:
		return "bloock_chain"
	case EthMainnet:
		return "ethereum_mainnet"
	case EthRinkeby:
		return "ethereum_rinkeby"
	case EthRopsten:
		return "ethereum_ropsten"
	case GnosisChain:
		return "gnosis_chain"
	case PolygonChain:
		return "polygon_chain"
	case EthGoerli:
		return "ethereum_goerli"
	case EthSepolia:
		return "ethereum_sepolia"
	default:
		return "ethereum_rinkeby"
	}
}
