package wrkchains

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/log"
	tmclient "github.com/tendermint/tendermint/rpc/client"
	tmtypes "github.com/tendermint/tendermint/types"
	"github.com/unification-com/wrkoracle/types"
)

// Tendermint is a structure for holding a Tendermint based WRKChain client
type Tendermint struct {
	log               log.Logger
	supportedHashMaps []string
}

// NewTendermintClient returns a new Tendermint struct
func NewTendermintClient() *Tendermint {
	return &Tendermint{
		supportedHashMaps: []string{"ReceiptHash", "TxHash", "Root", "UncleHash", "MixDigest"},
	}
}

// SetLogger sets the logger
func (t *Tendermint) SetLogger(log log.Logger) {
	t.log = log
}

// GetBlockAtHeight is used to get the block headers for a given height from a tendermint based WRKChain
func (t Tendermint) GetBlockAtHeight(height uint64) (WrkChainBlockHeader, error) {
	heightAt := int64(height)

	wrkChainClient, err := tmclient.NewHTTP(viper.GetString(types.FlagWrkchainRpc), "/websocket")
	if err != nil {
		return  WrkChainBlockHeader{}, err
	}

	if heightAt == 0 {
		status, err := wrkChainClient.Status()
		if err != nil {
			return WrkChainBlockHeader{}, err
		}
		heightAt = status.SyncInfo.LatestBlockHeight
	}

	latestWrkchainBlock, err := wrkChainClient.Block(&heightAt)

	if err != nil {
		return WrkChainBlockHeader{}, err
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
		hash1 = t.getHash(latestWrkchainBlock.Block.Header, hash1Ref)
	}

	if len(hash2Ref) > 0 {
		hash2 = t.getHash(latestWrkchainBlock.Block.Header, hash2Ref)
	}

	if len(hash3Ref) > 0 {
		hash3 = t.getHash(latestWrkchainBlock.Block.Header, hash3Ref)
	}

	wrkchainBlock := NewWrkChainBlockHeader(blockHeight, blockHash, parentHash, hash1, hash2, hash3)

	wrkChainClient.Quit()

	return wrkchainBlock, nil
}

// IsSupportedHash checks if the given hashType for the given chainType is currently supported by WRKOracle
func (t Tendermint) IsSupportedHash(hashType string) (bool, error) {
	for _, h := range t.supportedHashMaps {
		if hashType == h {
			return true, nil
		}
	}
	return false, fmt.Errorf("unsupported hash map '%s' for wrkchain type 'tendermint'. supported types: %s", hashType, strings.Join(t.supportedHashMaps, ", "))
}

// GetDefaultHashMap returns the default has mapping for a given reference
func (t Tendermint) GetDefaultHashMap(hashRef string) string {
	switch hashRef {
	case "hash1":
		return "DataHash"
	case "hash2":
		return "AppHash"
	case "hash3":
		return "ValidatorsHash"
	default:
		return ""
	}
}

func (t Tendermint) getHash(header tmtypes.Header, ref string) string {
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
		t.log.Error(fmt.Sprintf("unknown hash type '%s'", ref))
		return ""
	}
}
