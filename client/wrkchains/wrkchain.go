package wrkchains

import (
	"fmt"

	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/unification-com/wrkoracle/types"
)

var (
	// SupportedWrkchainTypes is an internal reference holding the currently supported WRKChain types
	SupportedWrkchainTypes = []string{
		"geth",
		"cosmos",
		"tendermint",
	}
)

// WrkChain is a top level struct to hold WRKChain data
type WrkChain struct {
	wrkChainClient WrkChainClient
	wrkchainMeta   WrkChainMeta
	log            log.Logger
}

// NewWrkChain returns a new initialised WrkChain
func NewWrkChain(wrkchainMeta WrkChainMeta, log log.Logger) (*WrkChain, error) {

	wrkChainClient, err := WrkChainClientFactory(wrkchainMeta.Type)
	if err != nil {
		return &WrkChain{}, err
	}

	wrkChainClient.SetLogger(log)

	return &WrkChain{
		wrkChainClient: wrkChainClient,
		wrkchainMeta:   wrkchainMeta,
		log:            log.With("pkg", "wrkchains"),
	}, nil
}

// WrkChainClientFactory returns a basic initialised WrkChainClient client
// based on the given WRKChain type
func WrkChainClientFactory(wrkchainType string) (WrkChainClient, error) {
	switch wrkchainType {
	case "geth":
		return NewGethClient(), nil
	case "tendermint", "cosmos":
		return NewTendermintClient(), nil
	default:
		var wrkChainClient WrkChainClient
		return wrkChainClient, fmt.Errorf("unsupported wrkchain type %s", wrkchainType)
	}
}

// IsSupportedHash checks if the given hashType for the given chainType is currently supported by WRKOracle
func IsSupportedHash(wrkchainType string, hashType string) (bool, error) {
	wrkChainClient, err := WrkChainClientFactory(wrkchainType)
	if err != nil {
		return false, err
	}
	return wrkChainClient.IsSupportedHash(hashType)
}

// GetDefaultHashMap returns the default hash map given a WRKChain type and hash reference (i.e. hash1, hash2 and hash3)
// It is called during the workoracle init command
func GetDefaultHashMap(wrkchainType string, hashRef string) string {
	wrkChainClient, err := WrkChainClientFactory(wrkchainType)
	if err != nil {
		return ""
	}
	return wrkChainClient.GetDefaultHashMap(hashRef)
}

// IsSupportedWrkchainType checks if the given chainType is currently supported by WRKOracle
func IsSupportedWrkchainType(wrkchainType string) bool {
	for _, wct := range SupportedWrkchainTypes {
		if wrkchainType == wct {
			return true
		}
	}
	return false
}

// GetLatestBlock is a top level function to query any WRKChain type for the latest block header
func (w WrkChain) GetLatestBlock() (WrkChainBlockHeader, error) {
	return w.GetWrkChainBlock(0)
}

// GetWrkChainBlock is a top level function to query any WRKChain type for the block header at a given height
func (w WrkChain) GetWrkChainBlock(height uint64) (WrkChainBlockHeader, error) {

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
