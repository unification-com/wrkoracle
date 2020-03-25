package wrkchains

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/unification-com/wrkoracle/types"
	"io/ioutil"
	"net/http"
	"strconv"
)

// nolint
const (
	TxMRoot     string = "TxMRoot"
	ActionRoot  string = "ActionRoot"
	ProducerSig string = "ProducerSig"
)

// EosGetBestBlockResult holds the result for a Eos getbestblockhash JSON RPC query
type EosGetBlockInfoResult struct {
	ChainId          string `json:"chain_id"`
	LastIrreversible uint64 `json:"last_irreversible_block_num"`
	HeadBlockNum     uint64 `json:"head_block_num"`
}

// EosBlockHeaderResult holds the result from a Eos JSON RPC query
type EosBlockHeaderResult struct {
	Id          string `json:"id"`        // hash
	BlockNum    uint64 `json:"block_num"` // height
	Previous    string `json:"previous"`  // parent hash
	TxMRoot     string `json:"transaction_mroot"`
	ActionRoot  string `json:"action_mroot"`
	ProducerSig string `json:"producer_signature"`
}

func init() {
	wrkchainClientCreator := func(log log.Logger, lastHeight uint64) WrkChainClient {
		return NewEosClient(log, lastHeight)
	}

	supportedHashMaps := []string{TxMRoot, ActionRoot, ProducerSig}

	defaultHashMap := make(map[string]string)
	defaultHashMap[types.FlagHash1] = TxMRoot
	defaultHashMap[types.FlagHash2] = ActionRoot
	defaultHashMap[types.FlagHash3] = ProducerSig

	registerWrkchainModule(EosWrkchainType, wrkchainClientCreator, supportedHashMaps, defaultHashMap, false)
}

var _ WrkChainClient = (*Eos)(nil)

// Eos is a structure for holding a Eos based WRKChain client
type Eos struct {
	log        log.Logger
	lastHeight uint64
}

// NewEosClient returns a new Eos struct
func NewEosClient(log log.Logger, lastHeight uint64) *Eos {
	return &Eos{
		log:        log,
		lastHeight: lastHeight,
	}
}

func (n Eos) getLatestBlockHash() (string, error) {
	queryUrl := viper.GetString(types.FlagWrkchainRpc) + "/v1/chain/get_info"
	resp, err := http.Post(queryUrl, "application/json", bytes.NewBuffer([]byte{}))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var res EosGetBlockInfoResult
	err = json.Unmarshal(body, &res)
	if err != nil {
		return "", err
	}
	return strconv.FormatUint(res.LastIrreversible, 10), nil
}

// GetBlockAtHeight is used to get the block headers for a given height from a Eos based WRKChain
func (n *Eos) GetBlockAtHeight(height uint64) (WrkChainBlockHeader, error) {

	queryUrl := viper.GetString(types.FlagWrkchainRpc) + "/v1/chain/get_block"

	var jsonStr []byte
	atHeight, err := n.getLatestBlockHash()
	if err != nil {
		return WrkChainBlockHeader{}, err
	}

	if height > 0 {
		atHeight = strconv.FormatUint(height, 10)
	}

	jsonStr = []byte(`{"block_num_or_id": "` + atHeight + `"}`)

	resp, err := http.Post(queryUrl, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		return WrkChainBlockHeader{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return WrkChainBlockHeader{}, err
	}

	var header EosBlockHeaderResult
	err = json.Unmarshal(body, &header)
	if err != nil {
		return WrkChainBlockHeader{}, err
	}

	if err != nil {
		return WrkChainBlockHeader{}, err
	}

	blockHash := header.Id
	parentHash := ""
	hash1 := ""
	hash2 := ""
	hash3 := ""
	blockHeight := header.BlockNum

	if err != nil {
		return WrkChainBlockHeader{}, err
	}

	if height == 0 {
		n.lastHeight = blockHeight
	}

	if viper.GetBool(types.FlagParentHash) {
		parentHash = header.Previous
	}

	hash1Ref := viper.GetString(types.FlagHash1)
	hash2Ref := viper.GetString(types.FlagHash2)
	hash3Ref := viper.GetString(types.FlagHash3)

	if len(hash1Ref) > 0 {
		hash1 = n.gethash(header, hash1Ref)
	}

	if len(hash2Ref) > 0 {
		hash2 = n.gethash(header, hash2Ref)
	}

	if len(hash3Ref) > 0 {
		hash3 = n.gethash(header, hash3Ref)
	}

	wrkchainBlock := NewWrkChainBlockHeader(blockHeight, blockHash, parentHash, hash1, hash2, hash3)

	return wrkchainBlock, nil
}

func (n Eos) gethash(header EosBlockHeaderResult, ref string) string {
	switch ref {
	case TxMRoot:
		return header.TxMRoot
	case ActionRoot:
		return header.ActionRoot
	case ProducerSig:
		return header.ProducerSig
	default:
		n.log.Error(fmt.Sprintf("unknown hash type '%s'", ref))
		return ""
	}
}
