package wrkchains

import (
	"github.com/unification-com/wrkoracle/types"
)

type WrkChainClient interface {
	GetBlockAtHeight(height uint64) (types.WrkChainBlockHeader, error)
}
