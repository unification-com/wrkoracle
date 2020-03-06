package geth

import (
	"context"
	"fmt"
	"math/big"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/viper"
	"github.com/unification-com/wrkoracle/types"
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
	hash1 := ""
	hash2 := ""
	hash3 := ""
	blockHeight := latestWrkchainHeader.Number.Uint64()

	if viper.GetBool(types.FlagParentHash) {
		parentHash = latestWrkchainHeader.ParentHash.String()
	}

	hash1Ref := viper.GetString(types.FlagHash1)
	hash2Ref := viper.GetString(types.FlagHash2)
	hash3Ref := viper.GetString(types.FlagHash3)

	if len(hash1Ref) > 0 {
		hash1 = getHash(latestWrkchainHeader, hash1Ref)
	}

	if len(hash2Ref) > 0 {
		hash2 = getHash(latestWrkchainHeader, hash2Ref)
	}

	if len(hash3Ref) > 0 {
		hash3 = getHash(latestWrkchainHeader, hash3Ref)
	}

	wrkchainBlock := types.NewWrkChainBlockHeader(blockHeight, blockHash, parentHash, hash1, hash2, hash3)

	return wrkchainBlock, nil
}

func getHash(header *ethtypes.Header, ref string) string {
	switch ref {
	case "ReceiptHash":
		return header.ReceiptHash.String()
	case "TxHash":
		return header.TxHash.String()
	case "Root":
		return header.Root.String()
	default:
		fmt.Println(fmt.Sprintf("unknown hash type '%s'", ref))
		return ""
	}
}
