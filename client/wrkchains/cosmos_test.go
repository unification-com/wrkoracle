package wrkchains

import (
	"github.com/spf13/viper"
	"github.com/unification-com/wrkoracle/types"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewCosmosClient(t *testing.T) {
	lh := uint64(1)
	cosCl := NewTendermintClient(logger, lh, CosmosWrkchainType)
	require.Equal(t, lh, cosCl.lastHeight)
	require.Equal(t, CosmosWrkchainType, cosCl.wrkchainType)
}

func TestCosmosGetWrkChainType(t *testing.T) {
	cosCl := NewTendermintClient(logger, 1, CosmosWrkchainType)
	wcT := cosCl.GetWrkChainType()
	require.Equal(t, CosmosWrkchainType, wcT)
}

func TestCosmosGetBlockAtHeight(t *testing.T) {
	h := uint64(1)
	cosCl := NewTendermintClient(logger, 1, CosmosWrkchainType)
	viper.Set(types.FlagWrkchainRpc, "http://18.218.225.121:26660")
	wcBlock, err := cosCl.GetBlockAtHeight(h)
	require.Nil(t, err)
	require.Equal(t, true, len(wcBlock.BlockHash) > 0)
	require.Equal(t, "B671921B3648F127D4F6900386083B6005885F947947647D3C4FA53811C1D405", wcBlock.BlockHash)
}
