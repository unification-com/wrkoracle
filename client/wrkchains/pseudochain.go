package wrkchains

import (
	"math/rand"
	"time"

	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/unification-com/wrkoracle/types"
)

// nolint
const (
	PseudoChainHash1 string = "PseudoChainHash1"
	PseudoChainHash2 string = "PseudoChainHash2"
	PseudoChainHash3 string = "PseudoChainHash3"
)

func init() {
	wrkchainClientCreator := func(log log.Logger, lastHeight uint64) WrkChainClient {
		return NewPseudoChainClient(log, lastHeight, PseudoChainWrkchainType)
	}

	supportedHashMaps := []string{PseudoChainHash1, PseudoChainHash2, PseudoChainHash3}

	defaultHashMap := make(map[string]string)
	defaultHashMap[types.FlagHash1] = PseudoChainHash1
	defaultHashMap[types.FlagHash2] = PseudoChainHash2
	defaultHashMap[types.FlagHash3] = PseudoChainHash3

	registerWrkchainModule(PseudoChainWrkchainType, wrkchainClientCreator, supportedHashMaps, defaultHashMap, false)
}

var _ WrkChainClient = (*PseudoChain)(nil)

// PseudoChain is a structure for holding a PseudoChain WRKChain client
type PseudoChain struct {
	log          log.Logger
	lastHeight   uint64
	parentHash   string
	seededRand   *rand.Rand
	wrkchainType WrkchainType
}

// NewPseudoChainClient returns a new PseudoChain struct
func NewPseudoChainClient(log log.Logger, lastHeight uint64, wrkchainType WrkchainType) *PseudoChain {
	return &PseudoChain{
		log:          log,
		lastHeight:   lastHeight,
		seededRand:   rand.New(rand.NewSource(time.Now().UnixNano())),
		wrkchainType: wrkchainType,
	}
}

// GetWrkChainType returns the WRKChain type
func (p PseudoChain) GetWrkChainType() WrkchainType {
	return p.wrkchainType
}

func (p PseudoChain) randomHash(length int) string {
	charset := "ABCDEF0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[p.seededRand.Intn(len(charset))]
	}
	return string(b)
}

// GetBlockAtHeight is used to get the block headers for a given height from a PseudoChain based WRKChain
func (p *PseudoChain) GetBlockAtHeight(height uint64) (WrkChainBlockHeader, error) {

	if height == 0 {
		height = p.lastHeight + 1
		p.lastHeight = height
	}

	blockHash := p.randomHash(64)
	parentHash := ""
	hash1 := ""
	hash2 := ""
	hash3 := ""
	blockHeight := height

	if viper.GetBool(types.FlagParentHash) {
		parentHash = p.parentHash
	}

	hash1Ref := viper.GetString(types.FlagHash1)
	hash2Ref := viper.GetString(types.FlagHash2)
	hash3Ref := viper.GetString(types.FlagHash3)

	if len(hash1Ref) > 0 {
		hash1 = p.randomHash(64)
	}

	if len(hash2Ref) > 0 {
		hash2 = p.randomHash(64)
	}

	if len(hash3Ref) > 0 {
		hash3 = p.randomHash(64)
	}

	wrkchainBlock := NewWrkChainBlockHeader(blockHeight, blockHash, parentHash, hash1, hash2, hash3)
	p.parentHash = blockHash

	return wrkchainBlock, nil
}
