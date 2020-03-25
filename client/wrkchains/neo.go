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
	MerkleRoot         string = "MerkleRoot"
	NextConsensus      string = "NextConsensus"
	NextBlockHash      string = "NextBlockHash"
	Nonce              string = "Nonce"
	ScriptInvocation   string = "ScriptInvocation"
	ScriptVerification string = "ScriptVerification"
)

type NeoGetBestBlockResult struct {
	Id      int    `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Result  string `json:"result"`
}

// NeoBlockHeaderResult holds the result from a Neo JSON RPC query
type NeoBlockHeaderResult struct {
	Id      int            `json:"id"`
	Jsonrpc string         `json:"jsonrpc"`
	Result  NeoBlockHeader `json:"result"`
}

// NeoBlockHeader holds the minimum Neo block header info returned from a Neo JSON RPC query
// required to process a Neo based WRKChain block header
type NeoBlockHeader struct {
	Index         uint64 `json:"index"`
	Hash          string `json:"hash"`
	PreviousHash  string `json:"previousblockhash"`
	MerkleRoot    string `json:"merkleroot"`
	NextConsensus string `json:"nextconsensus"`
	NextBlockHash string `json:"nextblockhash"`
	Nonce         string `json:"nonce"`
	Script        struct {
		Invocation   string `json:"invocation"`
		Verification string `json:"verification"`
	} `json:"script"`
}

func init() {
	wrkchainClientCreator := func(log log.Logger, lastHeight uint64) WrkChainClient {
		return NewNeoClient(log, lastHeight)
	}

	supportedHashMaps := []string{MerkleRoot, NextConsensus, NextBlockHash, Nonce, ScriptInvocation, ScriptVerification}

	defaultHashMap := make(map[string]string)
	defaultHashMap[types.FlagHash1] = MerkleRoot
	defaultHashMap[types.FlagHash2] = NextConsensus
	defaultHashMap[types.FlagHash3] = ScriptVerification

	registerWrkchainModule(NeoWrkchainType, wrkchainClientCreator, supportedHashMaps, defaultHashMap, false)
}

var _ WrkChainClient = (*Neo)(nil)

// Neo is a structure for holding a Neo based WRKChain client
type Neo struct {
	log        log.Logger
	lastHeight uint64
}

// NewNeoClient returns a new Neo struct
func NewNeoClient(log log.Logger, lastHeight uint64) *Neo {
	return &Neo{
		log:        log,
		lastHeight: lastHeight,
	}
}

func (n Neo) getLatestBlockHash() (string, error) {
	queryUrl := viper.GetString(types.FlagWrkchainRpc)
	var jsonStr = []byte(`{"jsonrpc":"2.0","method":"getbestblockhash","params":[], "id": 1}`)
	resp, err := http.Post(queryUrl, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var res NeoGetBestBlockResult
	err = json.Unmarshal(body, &res)
	if err != nil {
		return "", err
	}
	return res.Result, nil
}

// GetBlockAtHeight is used to get the block headers for a given height from a Neo based WRKChain
func (n *Neo) GetBlockAtHeight(height uint64) (WrkChainBlockHeader, error) {

	queryUrl := viper.GetString(types.FlagWrkchainRpc)

	var jsonStr []byte
	atHeight, err := n.getLatestBlockHash()
	if err != nil {
		return WrkChainBlockHeader{}, err
	}

	if height > 0 {
		atHeight = strconv.FormatUint(height, 10)
		jsonStr = []byte(`{"jsonrpc":"2.0","method":"getblock","params":[` + atHeight + `,1], "id": 1}`)
	} else {
		jsonStr = []byte(`{"jsonrpc":"2.0","method":"getblock","params":["` + atHeight + `",1], "id": 1}`)
	}

	resp, err := http.Post(queryUrl, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		return WrkChainBlockHeader{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return WrkChainBlockHeader{}, err
	}

	var res NeoBlockHeaderResult
	err = json.Unmarshal(body, &res)
	if err != nil {
		return WrkChainBlockHeader{}, err
	}

	header := res.Result

	if err != nil {
		return WrkChainBlockHeader{}, err
	}

	blockHash := header.Hash
	parentHash := ""
	hash1 := ""
	hash2 := ""
	hash3 := ""
	blockHeight := header.Index

	if height == 0 {
		n.lastHeight = blockHeight
	}

	if viper.GetBool(types.FlagParentHash) {
		parentHash = header.PreviousHash
	}

	hash1Ref := viper.GetString(types.FlagHash1)
	hash2Ref := viper.GetString(types.FlagHash2)
	hash3Ref := viper.GetString(types.FlagHash3)

	if len(hash1Ref) > 0 {
		hash1 = n.Neoash(header, hash1Ref)
	}

	if len(hash2Ref) > 0 {
		hash2 = n.Neoash(header, hash2Ref)
	}

	if len(hash3Ref) > 0 {
		hash3 = n.Neoash(header, hash3Ref)
	}

	wrkchainBlock := NewWrkChainBlockHeader(blockHeight, blockHash, parentHash, hash1, hash2, hash3)

	return wrkchainBlock, nil
}

func (n Neo) Neoash(header NeoBlockHeader, ref string) string {
	switch ref {
	case MerkleRoot:
		return header.MerkleRoot
	case NextConsensus:
		return header.NextConsensus
	case NextBlockHash:
		return header.NextBlockHash
	case Nonce:
		return header.Nonce
	case ScriptInvocation:
		return header.Script.Invocation
	case ScriptVerification:
		return header.Script.Verification
	default:
		n.log.Error(fmt.Sprintf("unknown hash type '%s'", ref))
		return ""
	}
}
