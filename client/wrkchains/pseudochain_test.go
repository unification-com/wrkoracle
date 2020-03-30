package wrkchains

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewPseudoChainClient(t *testing.T) {
	lh := uint64(1)
	pcC := NewPseudoChainClient(logger, lh, PseudoChainWrkchainType)
	require.Equal(t, lh, pcC.lastHeight)
	require.Equal(t, PseudoChainWrkchainType, pcC.wrkchainType)
}

func TestPseudoChainGetWrkChainType(t *testing.T) {
	pcC := NewPseudoChainClient(logger, 1, PseudoChainWrkchainType)
	wcT := pcC.GetWrkChainType()
	require.Equal(t, PseudoChainWrkchainType, wcT)
}

func TestPseudoChainGetBlockAtHeight(t *testing.T) {
	h := uint64(1)
	pcC := NewPseudoChainClient(logger, 1, PseudoChainWrkchainType)
	wcBlock, err := pcC.GetBlockAtHeight(h)
	require.Nil(t, err)

	require.Equal(t, true, len(wcBlock.BlockHash) > 0)
}
