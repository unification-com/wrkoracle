package wrkchains

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/unification-com/wrkoracle/types"
)

// WrkchainType is a string alias to hold supported WRKChain types
type WrkchainType string

// Currently supported WRKChain types
const (
	PseudoChainWrkchainType WrkchainType = "pseudochain" //generates random hashes for development only
	GethWrkchainType        WrkchainType = "geth"
	CosmosWrkchainType      WrkchainType = "cosmos"
	TendermintWrkchainType  WrkchainType = "tendermint"
	NeoWrkchainType         WrkchainType = "neo"
)

type wrkchainClientCreator func(log log.Logger, lastHeight uint64) WrkChainClient

// WrkchainModule is a structure to hold WRKChain module information
type WrkchainModule struct {
	creator        wrkchainClientCreator
	hashes         []string
	defaultHashMap map[string]string
}

// WrkChain is a top level struct to hold WRKChain data
type WrkChain struct {
	wrkChainClient WrkChainClient
	wrkchainMeta   WrkChainMeta
	log            log.Logger
}

var wrkchainModules = map[WrkchainType]WrkchainModule{}

// each module calls this to register it client creation method, supported optional hashes and default hash map in its init method.
func registerWrkchainModule(wrkchain WrkchainType, creator wrkchainClientCreator, hashes []string, hashMap map[string]string, force bool) {
	_, ok := wrkchainModules[wrkchain]
	if !force && ok {
		return
	}
	wrkchainModules[wrkchain] = WrkchainModule{creator: creator, hashes: hashes, defaultHashMap: hashMap}
}

// NewWrkChain returns a new initialised WrkChain
func NewWrkChain(wrkchainMeta WrkChainMeta, log log.Logger) (*WrkChain, error) {

	wrkChainModule, err := getWrkchainModule(wrkchainMeta.Type)

	if err != nil {
		return &WrkChain{}, err
	}

	lastHeight, err := strconv.Atoi(wrkchainMeta.LastBlock)
	if err != nil {
		lastHeight = 0
	}
	wrkChainClient := wrkChainModule.creator(log, uint64(lastHeight))

	return &WrkChain{
		wrkChainClient: wrkChainClient,
		wrkchainMeta:   wrkchainMeta,
		log:            log.With("pkg", "wrkchains"),
	}, nil
}

// IsSupportedHash checks if the given hashType for the given chainType is currently supported by WRKOracle
func IsSupportedHash(wrkchainType string, hashType string) (bool, error) {
	wrkChainModule, err := getWrkchainModule(wrkchainType)

	if err != nil {
		return false, err
	}

	for _, h := range wrkChainModule.hashes {
		if hashType == h {
			return true, nil
		}
	}
	return false, fmt.Errorf("unsupported hash map '%s' for wrkchain type '%s'. supported types: %s", hashType, wrkchainType, strings.Join(wrkChainModule.hashes, ", "))
}

// GetDefaultHashMap returns the default hash map given a WRKChain type and hash reference (i.e. hash1, hash2 and hash3)
// It is called during the workoracle init command
func GetDefaultHashMap(wrkchainType string, hashRef string) string {
	wrkChainModule, err := getWrkchainModule(wrkchainType)
	if err != nil {
		return ""
	}

	hash, ok := wrkChainModule.defaultHashMap[hashRef]
	if !ok {
		return ""
	}
	return hash
}

// IsSupportedWrkchainType checks if the given chainType is currently supported by WRKOracle
func IsSupportedWrkchainType(wrkchainType string) bool {
	_, ok := wrkchainModules[WrkchainType(wrkchainType)]
	return ok
}

// GetSupportedWrkchainTypes returns a slice of currently supported WRKChain types
func GetSupportedWrkchainTypes() []string {
	keys := make([]string, len(wrkchainModules))
	i := 0
	for k := range wrkchainModules {
		keys[i] = string(k)
		i++
	}
	return keys
}

// GetLatestBlock is a top level function to query any WRKChain type for the latest block header
func (w *WrkChain) GetLatestBlock() (WrkChainBlockHeader, error) {
	return w.GetWrkChainBlock(0)
}

// GetWrkChainBlock is a top level function to query any WRKChain type for the block header at a given height
func (w *WrkChain) GetWrkChainBlock(height uint64) (WrkChainBlockHeader, error) {

	w.log.Info("Get block for WRKChain", "moniker", w.wrkchainMeta.Moniker, "type", w.wrkchainMeta.Type, "rpc", viper.GetString(types.FlagWrkchainRpc))

	wrkchainBlock, err := w.wrkChainClient.GetBlockAtHeight(height)

	if err != nil {
		return WrkChainBlockHeader{}, err
	}

	w.log.Info("Got WRKChain block")
	w.log.Info("WRKChain Height", "height", wrkchainBlock.Height)
	w.log.Info("WRKChain Block Hash", "blockhash", wrkchainBlock.BlockHash)
	w.log.Info("WRKChain Parent Hash", "parenthash", wrkchainBlock.ParentHash)
	w.log.Info("WRKChain Hash1", "ref", viper.GetString(types.FlagHash1), "value", wrkchainBlock.Hash1)
	w.log.Info("WRKChain Hash2", "ref", viper.GetString(types.FlagHash2), "value", wrkchainBlock.Hash2)
	w.log.Info("WRKChain Hash3", "ref", viper.GetString(types.FlagHash3), "value", wrkchainBlock.Hash3)

	return wrkchainBlock, err
}

func getWrkchainModule(wrkchainType string) (WrkchainModule, error) {
	wrkChainModule, ok := wrkchainModules[WrkchainType(wrkchainType)]
	if !ok {
		keys := make([]string, len(wrkchainModules))
		i := 0
		for k := range wrkchainModules {
			keys[i] = string(k)
			i++
		}
		return WrkchainModule{}, fmt.Errorf("unknown wrkchain type %s, expected either %s", wrkchainType, strings.Join(keys, " or "))
	}
	return wrkChainModule, nil
}
