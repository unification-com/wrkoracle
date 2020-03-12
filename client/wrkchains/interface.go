package wrkchains

import (
	"github.com/tendermint/tendermint/libs/log"
)

// WrkChainClient is a generic interface for all WRKChain clients.
// New WRKChain client modules should implement this interface
type WrkChainClient interface {
	GetBlockAtHeight(height uint64) (WrkChainBlockHeader, error)
	IsSupportedHash(hashType string) (bool, error)
	GetDefaultHashMap(hashRef string) string
	SetLogger(log log.Logger)
}