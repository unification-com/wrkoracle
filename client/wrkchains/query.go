package wrkchains

import (
	"fmt"

	"github.com/spf13/viper"
	"github.com/unification-com/wrkoracle/client/wrkchains/geth"
	tm "github.com/unification-com/wrkoracle/client/wrkchains/tendermint"
	"github.com/unification-com/wrkoracle/types"
)

// GetLatestBlock is a top level function to query any WRKChain type for the latest block header
func GetLatestBlock(wrkchainMeta types.WrkChainMeta) (types.WrkChainBlockHeader, error) {
	return GetWrkChainBlock(wrkchainMeta, 0)
}

// GetWrkChainBlock is a top level function to query any WRKChain type for the block header at a given height
func GetWrkChainBlock(wrkchainMeta types.WrkChainMeta, height uint64) (types.WrkChainBlockHeader, error) {

	fmt.Println(fmt.Sprintf("Get block for WRKChain '%s', type '%s' at %s", wrkchainMeta.Moniker, wrkchainMeta.Type, viper.GetString(types.FlagWrkchainRpc)))

	var err error
	var wrkchainBlock types.WrkChainBlockHeader

	// generate a standard header object
	switch wrkchainMeta.Type {
	case "geth":
		wrkchainBlock, err = geth.GetBlock(height)
	case "tendermint", "cosmos":
		wrkchainBlock, err = tm.GetBlock(height)
	default:
		return types.WrkChainBlockHeader{}, fmt.Errorf("unsupported wrkchain type %s", wrkchainMeta.Type)
	}

	if err != nil {
		return types.WrkChainBlockHeader{}, err
	}

	hash1Ref := viper.GetString(types.FlagHash1)
	hash2Ref := viper.GetString(types.FlagHash2)
	hash3Ref := viper.GetString(types.FlagHash3)

	fmt.Println("Got WRKChain block")
	fmt.Println(fmt.Sprintf("WRKChain Height: %d", wrkchainBlock.Height))
	fmt.Println(fmt.Sprintf("WRKChain Block Hash: %s", wrkchainBlock.BlockHash))
	fmt.Println(fmt.Sprintf("WRKChain Parent Hash: %s", wrkchainBlock.ParentHash))
	fmt.Println(fmt.Sprintf("WRKChain Hash1 '%s': %s", hash1Ref, wrkchainBlock.Hash1))
	fmt.Println(fmt.Sprintf("WRKChain Hash2 '%s': %s", hash2Ref, wrkchainBlock.Hash2))
	fmt.Println(fmt.Sprintf("WRKChain Hash3 '%s': %s", hash3Ref, wrkchainBlock.Hash3))

	return wrkchainBlock, err
}
