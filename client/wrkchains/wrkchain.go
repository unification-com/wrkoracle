package wrkchains

import (
	"fmt"

	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/unification-com/wrkoracle/types"
)

// WrkChain is a top level struct to hold WRKChain data
type WrkChain struct {
	wrkChainClient WrkChainClient
	wrkchainMeta   types.WrkChainMeta
	log            log.Logger
}

// NewWrkChain returns a new initialised WrkChain
func NewWrkChain(wrkchainMeta types.WrkChainMeta, log log.Logger) (*WrkChain, error) {
	var wrkChainClient WrkChainClient

	switch wrkchainMeta.Type {
	case "geth":
		wrkChainClient = NewGethClient(log)
	case "tendermint", "cosmos":
		wrkChainClient = NewTendermintClient(log)
	default:
		return &WrkChain{}, fmt.Errorf("unsupported wrkchain type %s", wrkchainMeta.Type)
	}

	return &WrkChain{
		wrkChainClient: wrkChainClient,
		wrkchainMeta:   wrkchainMeta,
		log:            log.With("pkg", "wrkchains"),
	}, nil
}

// GetLatestBlock is a top level function to query any WRKChain type for the latest block header
func (w WrkChain) GetLatestBlock() (types.WrkChainBlockHeader, error) {
	return w.GetWrkChainBlock(0)
}

// GetWrkChainBlock is a top level function to query any WRKChain type for the block header at a given height
func (w WrkChain) GetWrkChainBlock(height uint64) (types.WrkChainBlockHeader, error) {

	w.log.Info("Get block for WRKChain", "moniker", w.wrkchainMeta.Moniker, "type", w.wrkchainMeta.Type, "rpc", viper.GetString(types.FlagWrkchainRpc))

	wrkchainBlock, err := w.wrkChainClient.GetBlockAtHeight(height)

	if err != nil {
		return types.WrkChainBlockHeader{}, err
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
