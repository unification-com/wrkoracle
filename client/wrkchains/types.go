package wrkchains

// WrkChainBlockHeader is the standard header object that should be returned by
// any WRKChain type client
type WrkChainBlockHeader struct {
	Height     uint64 `json:"height"`
	BlockHash  string `json:"blockhash"`
	ParentHash string `json:"parenthash"`
	Hash1      string `json:"hash1"`
	Hash2      string `json:"hash2"`
	Hash3      string `json:"hash3"`
}

// NewWrkChainBlockHeader returns a new initialised WrkChainBlockHeader
func NewWrkChainBlockHeader(
	height uint64,
	blockHash string,
	parentHash string,
	hash1 string,
	hash2 string,
	hash3 string,
) WrkChainBlockHeader {
	return WrkChainBlockHeader{
		Height:     height,
		BlockHash:  blockHash,
		ParentHash: parentHash,
		Hash1:      hash1,
		Hash2:      hash2,
		Hash3:      hash3,
	}
}

// WrkChainMeta is an object to hold WRKChain metadata when queried from Mainchain
type WrkChainMeta struct {
	WRKChainId string `json:"wrkchain_id"`
	Moniker    string `json:"moniker"`
	Name       string `json:"name"`
	Type       string `json:"type"`
}

// WrkChainMetaQueryResponse is a structure which holds Mainchain query data
// specifically for WRKChain metadata
type WrkChainMetaQueryResponse struct {
	Height string       `json:"height"`
	Result WrkChainMeta `json:"result"`
}
