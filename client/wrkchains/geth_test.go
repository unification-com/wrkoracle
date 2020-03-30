package wrkchains

import (
	"github.com/spf13/viper"
	"github.com/unification-com/wrkoracle/types"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewGethClient(t *testing.T) {
	lh := uint64(1)
	gethC := NewGethClient(logger, lh, GethWrkchainType)
	require.Equal(t, lh, gethC.lastHeight)
	require.Equal(t, GethWrkchainType, gethC.wrkchainType)
}

func TestGethGetWrkChainType(t *testing.T) {
	gethC := NewGethClient(logger, 1, GethWrkchainType)
	wcT := gethC.GetWrkChainType()
	require.Equal(t, GethWrkchainType, wcT)
}

func TestGethGetBlockAtHeight(t *testing.T) {
	h := uint64(7618305)
	gethC := NewGethClient(logger, 1, GethWrkchainType)
	viper.Set(types.FlagWrkchainRpc, "https://ropsten-rpc.linkpool.io")
	wcBlock, err := gethC.GetBlockAtHeight(h)
	require.Nil(t, err)

	require.Equal(t, true, len(wcBlock.BlockHash) > 0)
	require.Equal(t, "0x6bdc20cf17102ad7ae9eec5cd2041bc225b259c7c46c1a6160e5cc77348edca6", wcBlock.BlockHash)
}
