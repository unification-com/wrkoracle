package wrkchains

import (
	"encoding/json"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/unification-com/wrkoracle/types"
	"io/ioutil"
	"net/http"
	"strconv"
)

// nolint
const (
	HeaderXdr string = "HeaderXdr"
)

// StellarBlockHeader holds the minimum Stellar block header info returned from a Stellar JSON RPC query
// required to process a Stellar based WRKChain block header
type StellarBlockHeader struct {
	Hash      string `json:"hash"`
	PrevHash  string `json:"prev_hash"`
	Sequence  uint64 `json:"sequence"`
	HeaderXdr string `json:"header_xdr"`
}

// StellarLatestBlockHeader holds the minimum Stellar block header info returned from a Stellar "latest ledger" JSON RPC query
type StellarLatestBlockHeader struct {
	Embedded struct {
		Records []struct {
			StellarBlockHeader
		} `json:"records"`
	} `json:"_embedded"`
}

func init() {
	wrkchainClientCreator := func(log log.Logger, lastHeight uint64) WrkChainClient {
		return NewStellarClient(log, lastHeight)
	}

	supportedHashMaps := []string{HeaderXdr}

	defaultHashMap := make(map[string]string)
	defaultHashMap[types.FlagHash1] = HeaderXdr
	defaultHashMap[types.FlagHash2] = ""
	defaultHashMap[types.FlagHash3] = ""

	registerWrkchainModule(StellarWrkchainType, wrkchainClientCreator, supportedHashMaps, defaultHashMap, false)
}

var _ WrkChainClient = (*Stellar)(nil)

// Stellar is a structure for holding a Stellar based WRKChain client
type Stellar struct {
	log        log.Logger
	lastHeight uint64
}

// NewStellarClient returns a new Stellar struct
func NewStellarClient(log log.Logger, lastHeight uint64) *Stellar {
	return &Stellar{
		log:        log,
		lastHeight: lastHeight,
	}
}

// GetBlockAtHeight is used to get the block headers for a given height from a Stellar based WRKChain
func (n *Stellar) GetBlockAtHeight(height uint64) (WrkChainBlockHeader, error) {

	queryUrl := viper.GetString(types.FlagWrkchainRpc) + "/ledgers"

	if height > 0 {
		queryUrl = queryUrl + "/" + strconv.Itoa(int(height))
	} else {
		queryUrl = queryUrl + "?cursor=now&order=desc&limit=1"
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

	hash1Ref := viper.GetString(types.FlagHash1)

	blockHash := ""
	parentHash := ""
	hash1 := ""
	hash2 := ""
	hash3 := ""
	blockHeight := uint64(0)

	if height == 0 {
		var res StellarLatestBlockHeader
		err = json.Unmarshal(body, &res)
		if err != nil {
			return WrkChainBlockHeader{}, err
		}
		blockHash = res.Embedded.Records[0].Hash
		blockHeight = res.Embedded.Records[0].Sequence
		if viper.GetBool(types.FlagParentHash) {
			parentHash = res.Embedded.Records[0].PrevHash
		}
		if len(hash1Ref) > 0 {
			hash1 = res.Embedded.Records[0].HeaderXdr
		}
		n.lastHeight = blockHeight
	} else {
		var res StellarBlockHeader
		err = json.Unmarshal(body, &res)
		if err != nil {
			return WrkChainBlockHeader{}, err
		}
		blockHash = res.Hash
		blockHeight = res.Sequence
		if viper.GetBool(types.FlagParentHash) {
			parentHash = res.PrevHash
		}
		if len(hash1Ref) > 0 {
			hash1 = res.HeaderXdr
		}
	}

	wrkchainBlock := NewWrkChainBlockHeader(blockHeight, blockHash, parentHash, hash1, hash2, hash3)

	return wrkchainBlock, nil
}
