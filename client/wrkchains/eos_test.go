package wrkchains

import (
	"github.com/spf13/viper"
	"github.com/unification-com/wrkoracle/types"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewEosClient(t *testing.T) {
	lh := uint64(1)
	eosC := NewEosClient(logger, lh, EosWrkchainType)
	require.Equal(t, lh, eosC.lastHeight)
	require.Equal(t, EosWrkchainType, eosC.wrkchainType)
}

func TestEosGetWrkChainType(t *testing.T) {
	eosC := NewEosClient(logger, 1, EosWrkchainType)
	wcT := eosC.GetWrkChainType()
	require.Equal(t, EosWrkchainType, wcT)
}

func TestEosGetBlockAtHeight(t *testing.T) {
	h := uint64(1)
	eosC := NewEosClient(logger, 1, EosWrkchainType)
	viper.Set(types.FlagWrkchainRpc, "https://api.testnet.eos.io")
	wcBlock, err := eosC.GetBlockAtHeight(h)
	require.Nil(t, err)
	require.Equal(t, true, len(wcBlock.BlockHash) > 0)
	require.Equal(t, "00000001c445c471d0b5563b0f73a113b14fcd3409df821c4ed42d32c0e3c5ca", wcBlock.BlockHash)
}
