package wrkchains

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/unification-com/wrkoracle/types"
)

// nolint
const (
	ReceiptsRoot string = "ReceiptsRoot"
	TxRoot       string = "TxRoot"
	StateRoot    string = "StateRoot"
	UncleHash    string = "UncleHash"
	MixHash      string = "MixHash"
)

// GethBlockHeaderResult holds the result from a Geth JSON RPC query
type GethBlockHeaderResult struct {
	Id      string          `json:"id"`
	Jsonrpc string          `json:"jsonrpc"`
	Result  GethBlockHeader `json:"result"`
}

// GethBlockHeader holds the minimum Geth block header info returned from a Geth JSON RPC query
// required to process a geth based WRKChain block header
type GethBlockHeader struct {
	Number       string `json:"number"`
	Hash         string `json:"hash"`
	ParentHash   string `json:"parentHash"`
	MixHash      string `json:"mixHash"`
	UncleHash    string `json:"sha3Uncles"`
	TxRoot       string `json:"transactionsRoot"`
	StateRoot    string `json:"stateRoot"`
	ReceiptsRoot string `json:"receiptsRoot"`
}

func init() {
	wrkchainClientCreator := func(log log.Logger, lastHeight uint64) WrkChainClient {
		return NewGethClient(log, lastHeight)
	}

	supportedHashMaps := []string{ReceiptsRoot, TxRoot, StateRoot, UncleHash, MixHash}

	defaultHashMap := make(map[string]string)
	defaultHashMap[types.FlagHash1] = ReceiptsRoot
	defaultHashMap[types.FlagHash2] = TxRoot
	defaultHashMap[types.FlagHash3] = StateRoot

	registerWrkchainModule(GethWrkchainType, wrkchainClientCreator, supportedHashMaps, defaultHashMap, false)
}

var _ WrkChainClient = (*Geth)(nil)

// Geth is a structure for holding a Geth based WRKChain client
type Geth struct {
	log        log.Logger
	lastHeight uint64
}

// NewGethClient returns a new Geth struct
func NewGethClient(log log.Logger, lastHeight uint64) *Geth {
	return &Geth{
		log:        log,
		lastHeight: lastHeight,
	}
}

// GetBlockAtHeight is used to get the block headers for a given height from a geth based WRKChain
func (g *Geth) GetBlockAtHeight(height uint64) (WrkChainBlockHeader, error) {

	queryUrl := viper.GetString(types.FlagWrkchainRpc)

	atHeight := "latest"

	if height > 0 {
		atHeight = "0x" + strconv.FormatUint(height, 16)
	}

	var jsonStr = []byte(`{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["` + atHeight + `",false]}`)

	resp, err := http.Post(queryUrl, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		return WrkChainBlockHeader{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return WrkChainBlockHeader{}, err
	}

	var res GethBlockHeaderResult
	err = json.Unmarshal(body, &res)
	if err != nil {
		return WrkChainBlockHeader{}, err
	}

	header := res.Result
	cleanedHeight := strings.Replace(header.Number, "0x", "", -1)
	blockNumber, err := strconv.ParseUint(cleanedHeight, 16, 64)

	if err != nil {
		return WrkChainBlockHeader{}, err
	}

	blockHash := header.Hash
	parentHash := ""
	hash1 := ""
	hash2 := ""
	hash3 := ""
	blockHeight := blockNumber

	if height == 0 {
		g.lastHeight = blockNumber
	}

	if viper.GetBool(types.FlagParentHash) {
		parentHash = header.ParentHash
	}

	hash1Ref := viper.GetString(types.FlagHash1)
	hash2Ref := viper.GetString(types.FlagHash2)
	hash3Ref := viper.GetString(types.FlagHash3)

	if len(hash1Ref) > 0 {
		hash1 = g.getHash(header, hash1Ref)
	}

	if len(hash2Ref) > 0 {
		hash2 = g.getHash(header, hash2Ref)
	}

	if len(hash3Ref) > 0 {
		hash3 = g.getHash(header, hash3Ref)
	}

	wrkchainBlock := NewWrkChainBlockHeader(blockHeight, blockHash, parentHash, hash1, hash2, hash3)

	return wrkchainBlock, nil
}

func (g Geth) getHash(header GethBlockHeader, ref string) string {
	switch ref {
	case ReceiptsRoot:
		return header.ReceiptsRoot
	case TxRoot:
		return header.TxRoot
	case StateRoot:
		return header.StateRoot
	case UncleHash:
		return header.UncleHash
	case MixHash:
		return header.MixHash
	default:
		g.log.Error(fmt.Sprintf("unknown hash type '%s'", ref))
		return ""
	}
}
