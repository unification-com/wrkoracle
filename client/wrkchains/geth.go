package wrkchains

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/unification-com/wrkoracle/types"
)

// Geth is a structure for holding a Geth based WRKChain client
type Geth struct {
	log               log.Logger
	supportedHashMaps []string
}

// NewGethClient returns a new Geth struct
func NewGethClient() *Geth {
	return &Geth{
		supportedHashMaps: []string{"ReceiptHash", "TxHash", "Root", "UncleHash", "MixDigest"},
	}
}

// SetLogger sets the logger
func (g *Geth) SetLogger(log log.Logger) {
	g.log = log
}

// GetBlockAtHeight is used to get the block headers for a given height from a geth based WRKChain
func (g Geth) GetBlockAtHeight(height uint64) (WrkChainBlockHeader, error) {

	wrkChainClient, _ := ethclient.Dial(viper.GetString(types.FlagWrkchainRpc))

	atHeight := big.NewInt(0).SetUint64(height)

	if height == 0 {
		atHeight = nil
	}

	latestWrkchainHeader, err := wrkChainClient.HeaderByNumber(context.Background(), atHeight)

	if err != nil {
		return WrkChainBlockHeader{}, err
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
		hash1 = g.getHash(latestWrkchainHeader, hash1Ref)
	}

	if len(hash2Ref) > 0 {
		hash2 = g.getHash(latestWrkchainHeader, hash2Ref)
	}

	if len(hash3Ref) > 0 {
		hash3 = g.getHash(latestWrkchainHeader, hash3Ref)
	}

	wrkchainBlock := NewWrkChainBlockHeader(blockHeight, blockHash, parentHash, hash1, hash2, hash3)

	return wrkchainBlock, nil
}

// IsSupportedHash checks if the given hashType for the given chainType is currently supported by WRKOracle
func (g Geth) IsSupportedHash(hashType string) (bool, error) {
	for _, h := range g.supportedHashMaps {
		if hashType == h {
			return true, nil
		}
	}
	return false, fmt.Errorf("unsupported hash map '%s' for wrkchain type 'geth'. supported types: %s", hashType, strings.Join(g.supportedHashMaps, ", "))
}

// GetDefaultHashMap returns the default has mapping for a given reference
func (g Geth) GetDefaultHashMap(hashRef string) string {
	switch hashRef {
	case "hash1":
		return "ReceiptHash"
	case "hash2":
		return "TxHash"
	case "hash3":
		return "Root"
	default:
		return ""
	}
}

func (g Geth) getHash(header *ethtypes.Header, ref string) string {
	switch ref {
	case "ReceiptHash":
		return header.ReceiptHash.String()
	case "TxHash":
		return header.TxHash.String()
	case "Root":
		return header.Root.String()
	case "UncleHash":
		return header.UncleHash.String()
	case "MixDigest":
		return header.MixDigest.String()
	default:
		g.log.Error(fmt.Sprintf("unknown hash type '%s'", ref))
		return ""
	}
}
