package geth

import (
	"context"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/viper"
	"github.com/unification-com/wrkoracle/types"
	"math/big"
)

// GetBlock is used to get the block headers for a given height from a geth based WRKChain
func GetBlock(height uint64) (types.WrkChainBlockHeader, error) {

	wrkChainClient, _ := ethclient.Dial(viper.GetString(types.FlagWrkchainRpc))

	atHeight := big.NewInt(0).SetUint64(height)

	if height == 0 {
		atHeight = nil
	}

	latestWrkchainHeader, err := wrkChainClient.HeaderByNumber(context.Background(), atHeight)

	if err != nil {
		return types.WrkChainBlockHeader{}, err
	}

	blockHash := latestWrkchainHeader.Hash().String()
	parentHash := ""
	receiptHash := ""
	txHash := ""
	rootHash := ""
	blockHeight := latestWrkchainHeader.Number.Uint64()

	if viper.GetBool(types.FlagParentHash) {
		parentHash = latestWrkchainHeader.ParentHash.String()
	}

	if viper.GetBool(types.FlagHash1) {
		receiptHash = latestWrkchainHeader.ReceiptHash.String()
	}

	if viper.GetBool(types.FlagHash2) {
		txHash = latestWrkchainHeader.TxHash.String()
	}

	if viper.GetBool(types.FlagHash3) {
		rootHash = latestWrkchainHeader.Root.String()
	}

	wrkchainBlock := types.NewWrkChainBlockHeader(blockHeight, blockHash, parentHash, receiptHash, txHash, rootHash)

	return wrkchainBlock, nil
}
