package wrkchains

import (
	"github.com/spf13/viper"
	"github.com/unification-com/wrkoracle/types"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewStellarClient(t *testing.T) {
	lh := uint64(1)
	stellarC := NewStellarClient(logger, lh, StellarWrkchainType)
	require.Equal(t, lh, stellarC.lastHeight)
	require.Equal(t, StellarWrkchainType, stellarC.wrkchainType)
}

func TestStellarGetWrkChainType(t *testing.T) {
	stellarC := NewStellarClient(logger, 1, StellarWrkchainType)
	wcT := stellarC.GetWrkChainType()
	require.Equal(t, StellarWrkchainType, wcT)
}

func TestStellarGetBlockAtHeight(t *testing.T) {
	h := uint64(1)
	stellarC := NewStellarClient(logger, 1, StellarWrkchainType)
	viper.Set(types.FlagWrkchainRpc, "https://horizon-testnet.stellar.org")
	wcBlock, err := stellarC.GetBlockAtHeight(h)
	require.Nil(t, err)

	require.Equal(t, true, len(wcBlock.BlockHash) > 0)
	require.Equal(t, "63d98f536ee68d1b27b5b89f23af5311b7569a24faf1403ad0b52b633b07be99", wcBlock.BlockHash)
}
