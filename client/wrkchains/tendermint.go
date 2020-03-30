package wrkchains

import (
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
	DataHash           string = "DataHash"
	AppHash            string = "AppHash"
	ValidatorsHash     string = "ValidatorsHash"
	LastResultsHash    string = "LastResultsHash"
	LastCommitHash     string = "LastCommitHash"
	ConsensusHash      string = "ConsensusHash"
	NextValidatorsHash string = "NextValidatorsHash"
	EvidenceHash       string = "EvidenceHash"
)

// TmBlockHeaderResult holds the result from a Tendermint node RPC query
type TmBlockHeaderResult struct {
	Result TmResult `json:"result"`
}

// TmResult holds the minimum amount of data returned from a Tendermint node RPC request
type TmResult struct {
	BlockID TmBlockID `json:"block_id"`
	Block   TmBlock   `json:"block"`
}

// TmBlockID holds the minimum amount of block ID data returned from a Tendermint node RPC request
type TmBlockID struct {
	Hash string `json:"hash"`
}

// TmBlock holds the minimum amount of block ID data returned from a Tendermint node RPC request
type TmBlock struct {
	Header TmBlockHeader `json:"header"`
}

// TmBlockHeader holds the minimum Tendermint block header info returned from a TM RPC query
// required to process a geth based WRKChain block header
type TmBlockHeader struct {
	// prev block info
	LastBlockId TmBlockID `json:"last_block_id"`

	Height  string `json:"height"`
	ChainId string `json:"chain_id"`

	// hashes of block data
	LastCommitHash string `json:"last_commit_hash"` // commit from validators from the last block
	DataHash       string `json:"data_hash"`        // transactions

	// hashes from the app output from the prev block
	ValidatorsHash     string `json:"validators_hash"`      // validators for the current block
	NextValidatorsHash string `json:"next_validators_hash"` // validators for the next block
	ConsensusHash      string `json:"consensus_hash"`       // consensus params for current block
	AppHash            string `json:"app_hash"`             // state after txs from the previous block
	// root hash of all results from the txs from the previous block
	LastResultsHash string `json:"last_results_hash"`
	// consensus info
	EvidenceHash string `json:"evidence_hash"` // evidence included in the block
}

func init() {
	wrkchainClientCreator := func(log log.Logger, lastHeight uint64) WrkChainClient {
		return NewTendermintClient(log, lastHeight, TendermintWrkchainType)
	}

	supportedHashMaps := []string{DataHash, AppHash, ValidatorsHash, LastResultsHash, LastCommitHash, ConsensusHash, NextValidatorsHash, EvidenceHash}

	defaultHashMap := make(map[string]string)
	defaultHashMap[types.FlagHash1] = DataHash
	defaultHashMap[types.FlagHash2] = AppHash
	defaultHashMap[types.FlagHash3] = ValidatorsHash

	registerWrkchainModule(TendermintWrkchainType, wrkchainClientCreator, supportedHashMaps, defaultHashMap, false)
}

var _ WrkChainClient = (*Tendermint)(nil)

// Tendermint is a structure for holding a Tendermint based WRKChain client
type Tendermint struct {
	log          log.Logger
	lastHeight   uint64
	wrkchainType WrkchainType
}

// NewTendermintClient returns a new Tendermint struct
func NewTendermintClient(log log.Logger, lastHeight uint64, wrkchainType WrkchainType) *Tendermint {
	return &Tendermint{
		log:          log,
		lastHeight:   lastHeight,
		wrkchainType: wrkchainType,
	}
}

// GetWrkChainType returns the WRKChain type
func (t Tendermint) GetWrkChainType() WrkchainType {
	return t.wrkchainType
}

// GetBlockAtHeight is used to get the block headers for a given height from a tendermint based WRKChain
func (t *Tendermint) GetBlockAtHeight(height uint64) (WrkChainBlockHeader, error) {

	queryUrl := viper.GetString(types.FlagWrkchainRpc) + "/block"
	if height > 0 {
		queryUrl = queryUrl + "?height=" + strconv.Itoa(int(height))
	}

	resp, err := http.Get(queryUrl)
	if err != nil {
		return WrkChainBlockHeader{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return WrkChainBlockHeader{}, err
	}

	var res TmBlockHeaderResult

	err = json.Unmarshal(body, &res)

	if err != nil {
		return WrkChainBlockHeader{}, err
	}

	tmBlock := res.Result

	blockHash := tmBlock.BlockID.Hash

	parentHash := ""
	hash1 := ""
	hash2 := ""
	hash3 := ""
	blockHeight, err := strconv.Atoi(tmBlock.Block.Header.Height)

	if height == 0 {
		t.lastHeight = uint64(blockHeight)
	}

	if err != nil {
		return WrkChainBlockHeader{}, err
	}

	if viper.GetBool(types.FlagParentHash) {
		parentHash = tmBlock.Block.Header.LastBlockId.Hash
	}

	hash1Ref := viper.GetString(types.FlagHash1)
	hash2Ref := viper.GetString(types.FlagHash2)
	hash3Ref := viper.GetString(types.FlagHash3)

	if len(hash1Ref) > 0 {
		hash1 = t.getHash(tmBlock.Block.Header, hash1Ref)
	}

	if len(hash2Ref) > 0 {
		hash2 = t.getHash(tmBlock.Block.Header, hash2Ref)
	}

	if len(hash3Ref) > 0 {
		hash3 = t.getHash(tmBlock.Block.Header, hash3Ref)
	}

	wrkchainBlock := NewWrkChainBlockHeader(uint64(blockHeight), blockHash, parentHash, hash1, hash2, hash3)

	return wrkchainBlock, nil
}

func (t Tendermint) getHash(header TmBlockHeader, ref string) string {
	switch ref {
	case DataHash:
		return header.DataHash
	case AppHash:
		return header.AppHash
	case ValidatorsHash:
		return header.ValidatorsHash
	case LastResultsHash:
		return header.LastResultsHash
	case LastCommitHash:
		return header.LastCommitHash
	case ConsensusHash:
		return header.ConsensusHash
	case NextValidatorsHash:
		return header.NextValidatorsHash
	default:
		t.log.Error(fmt.Sprintf("unknown hash type '%s'", ref))
		return ""
	}
}
