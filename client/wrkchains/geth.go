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

// GethBlockHeaderResult holds the result from a Geth JSON RPC query
type GethBlockHeaderResult struct {
    Id string `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
    Result GethBlockHeader  `json:"result"`
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

// Geth is a structure for holding a Geth based WRKChain client
type Geth struct {
	log               log.Logger
	supportedHashMaps []string
}

// NewGethClient returns a new Geth struct
func NewGethClient() *Geth {
	return &Geth{
		supportedHashMaps: []string{"ReceiptsRoot", "TxRoot", "StateRoot", "UncleHash", "MixHash"},
	}
}

// SetLogger sets the logger
func (g *Geth) SetLogger(log log.Logger) {
	g.log = log
}

// GetBlockAtHeight is used to get the block headers for a given height from a geth based WRKChain
func (g Geth) GetBlockAtHeight(height uint64) (WrkChainBlockHeader, error) {

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
		return WrkChainBlockHeader{}, nil
	}

	header := res.Result
	cleanedHeight := strings.Replace(header.Number, "0x", "", -1)
	blockNumber, err := strconv.ParseUint(cleanedHeight, 16, 64)

	if err != nil {
		return WrkChainBlockHeader{}, nil
	}

	blockHash := header.Hash
	parentHash := ""
	hash1 := ""
	hash2 := ""
	hash3 := ""
	blockHeight := blockNumber

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

// IsSupportedHash checks if the given hashType for the given chainType is currently supported by WRKOracle
func (g Geth) IsSupportedHash(hashType string) (bool, error) {
	for _, h := range g.supportedHashMaps {
		if hashType == h {
			return true, nil
		}
	}
	return false, fmt.Errorf("unsupported hash map '%s' for wrkchain type 'geth'. supported types: %s", hashType, strings.Join(g.supportedHashMaps, ", "))
}

// GetDefaultHashMap returns the default has mapping for a given reference
func (g Geth) GetDefaultHashMap(hashRef string) string {
	switch hashRef {
	case "hash1":
		return "ReceiptsRoot"
	case "hash2":
		return "TxRoot"
	case "hash3":
		return "StateRoot"
	default:
		return ""
	}
}

func (g Geth) getHash(header GethBlockHeader, ref string) string {
	switch ref {
	case "ReceiptsRoot":
		return header.ReceiptsRoot
	case "TxRoot":
		return header.TxRoot
	case "StateRoot":
		return header.StateRoot
	case "UncleHash":
		return header.UncleHash
	case "MixHash":
		return header.MixHash
	default:
		g.log.Error(fmt.Sprintf("unknown hash type '%s'", ref))
		return ""
	}
}
