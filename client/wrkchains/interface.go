package wrkchains

import (
	"github.com/unification-com/wrkoracle/types"
)

// WrkChainClient is a generic interface for all WRKChain clients.
// New WRKChain client modules should implement this interface
type WrkChainClient interface {
	GetBlockAtHeight(height uint64) (types.WrkChainBlockHeader, error)
}
