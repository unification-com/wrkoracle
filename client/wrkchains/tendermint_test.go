package wrkchains

import (
	"github.com/spf13/viper"
	"github.com/unification-com/wrkoracle/types"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewTendermintClient(t *testing.T) {
	lh := uint64(1)
	tmC := NewTendermintClient(logger, lh, TendermintWrkchainType)
	require.Equal(t, lh, tmC.lastHeight)
	require.Equal(t, TendermintWrkchainType, tmC.wrkchainType)
}

func TestTendermintGetWrkChainType(t *testing.T) {
	tmC := NewTendermintClient(logger, 1, TendermintWrkchainType)
	wcT := tmC.GetWrkChainType()
	require.Equal(t, TendermintWrkchainType, wcT)
}

func TestTendermintGetBlockAtHeight(t *testing.T) {
	h := uint64(1)
	tmC := NewTendermintClient(logger, 1, TendermintWrkchainType)
	viper.Set(types.FlagWrkchainRpc, "http://18.218.225.121:26660")
	wcBlock, err := tmC.GetBlockAtHeight(h)
	require.Nil(t, err)
	require.Equal(t, true, len(wcBlock.BlockHash) > 0)
	require.Equal(t, "B671921B3648F127D4F6900386083B6005885F947947647D3C4FA53811C1D405", wcBlock.BlockHash)
}
