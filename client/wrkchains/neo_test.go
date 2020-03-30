package wrkchains

import (
	"github.com/spf13/viper"
	"github.com/unification-com/wrkoracle/types"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewNeoClient(t *testing.T) {
	lh := uint64(1)
	neoC := NewNeoClient(logger, lh, NeoWrkchainType)
	require.Equal(t, lh, neoC.lastHeight)
	require.Equal(t, NeoWrkchainType, neoC.wrkchainType)
}

func TestNeoGetWrkChainType(t *testing.T) {
	neoC := NewNeoClient(logger, 1, NeoWrkchainType)
	wcT := neoC.GetWrkChainType()
	require.Equal(t, NeoWrkchainType, wcT)
}

func TestNeoGetBlockAtHeight(t *testing.T) {
	h := uint64(1)
	neoC := NewNeoClient(logger, 1, NeoWrkchainType)
	viper.Set(types.FlagWrkchainRpc, "http://seed2.ngd.network:10332")
	wcBlock, err := neoC.GetBlockAtHeight(h)
	require.Nil(t, err)

	require.Equal(t, true, len(wcBlock.BlockHash) > 0)
	require.Equal(t, "0xd782db8a38b0eea0d7394e0f007c61c71798867578c77c387c08113903946cc9", wcBlock.BlockHash)
}
