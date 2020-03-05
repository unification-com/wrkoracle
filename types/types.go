package types

type WrkChainBlockHeader struct {
	Height     uint64 `json:"height"`
	BlockHash  string `json:"blockhash"`
	ParentHash string `json:"parenthash"`
	Hash1      string `json:"hash1"`
	Hash2      string `json:"hash2"`
	Hash3      string `json:"hash3"`
}

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

type WrkChainMeta struct {
	WRKChainId string `json:"wrkchain_id"`
	Moniker    string `json:"moniker"`
	Name       string `json:"name"`
	Type       string `json:"type"`
}

type WrkChainMetaQueryResponse struct {
	Height string       `json:"height"`
	Result WrkChainMeta `json:"result"`
}

type FeeParams struct {
	Denom       string `json:"denom"`
	FeeRecord   string `json:"fee_record"`
	FeeRegister string `json:"fee_register"`
}

type FeeParamsQueryResponse struct {
	Height string    `json:"height"`
	Result FeeParams `json:"result"`
}
