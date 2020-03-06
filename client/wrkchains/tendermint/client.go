package tendermint

import (
	"fmt"

	"github.com/spf13/viper"
	tmclient "github.com/tendermint/tendermint/rpc/client"
	tmtypes "github.com/tendermint/tendermint/types"
	"github.com/unification-com/wrkoracle/types"
)

// GetBlock is used to get the block headers for a given height from a tendermint based WRKChain
func GetBlock(height uint64) (types.WrkChainBlockHeader, error) {
	heightAt := int64(height)
	wrkChainClient, err := tmclient.NewHTTP(viper.GetString(types.FlagWrkchainRpc), "/websocket")

	if err != nil {
		return types.WrkChainBlockHeader{}, err
	}

	if heightAt == 0 {
		status, err := wrkChainClient.Status()
		if err != nil {
			return types.WrkChainBlockHeader{}, err
		}
		heightAt = status.SyncInfo.LatestBlockHeight
	}

	latestWrkchainBlock, err := wrkChainClient.Block(&heightAt)

	if err != nil {
		return types.WrkChainBlockHeader{}, err
	}

	blockHash := latestWrkchainBlock.BlockID.Hash.String()

	parentHash := ""
	hash1 := ""
	hash2 := ""
	hash3 := ""
	blockHeight := uint64(latestWrkchainBlock.Block.Height)

	if viper.GetBool(types.FlagParentHash) {
		parentHash = latestWrkchainBlock.Block.Header.LastBlockID.Hash.String()
	}

	hash1Ref := viper.GetString(types.FlagHash1)
	hash2Ref := viper.GetString(types.FlagHash2)
	hash3Ref := viper.GetString(types.FlagHash3)

	if len(hash1Ref) > 0 {
		hash1 = getHash(latestWrkchainBlock.Block.Header, hash1Ref)
	}

	if len(hash2Ref) > 0 {
		hash2 = getHash(latestWrkchainBlock.Block.Header, hash2Ref)
	}

	if len(hash3Ref) > 0 {
		hash3 = getHash(latestWrkchainBlock.Block.Header, hash3Ref)
	}

	wrkchainBlock := types.NewWrkChainBlockHeader(blockHeight, blockHash, parentHash, hash1, hash2, hash3)

	return wrkchainBlock, nil
}

func getHash(header tmtypes.Header, ref string) string {
	switch ref {
	case "DataHash":
		return header.DataHash.String()
	case "AppHash":
		return header.AppHash.String()
	case "ValidatorsHash":
		return header.ValidatorsHash.String()
	case "LastResultsHash":
		return header.LastResultsHash.String()
	case "LastCommitHash":
		return header.LastCommitHash.String()
	case "ConsensusHash":
		return header.ConsensusHash.String()
	case "NextValidatorsHash":
		return header.NextValidatorsHash.String()
	default:
		fmt.Println(fmt.Sprintf("unknown hash type '%s'", ref))
		return ""
	}
}
