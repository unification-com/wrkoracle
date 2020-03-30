package wrkchains

import (
	"fmt"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"github.com/unification-com/wrkoracle/types"
)

func TestNewWrkChain(t *testing.T) {

	wcTypes := GetSupportedWrkchainTypes()
	var wrkchainMeta = WrkChainMeta{
		LastBlock: "1",
	}
	for _, wt := range wcTypes {
		wrkchainMeta.Type = wt

		wc, err := NewWrkChain(wrkchainMeta, logger)

		require.Nil(t, err)

		require.Equal(t, wt, string(wc.WrkChainClient.GetWrkChainType()))
	}

	wrkchainMeta.Type = "rubbishchain"
	_, err := NewWrkChain(wrkchainMeta, logger)

	msg := fmt.Errorf("unknown wrkchain type rubbishchain, expected either %s", strings.Join(wcTypes, " or "))
	require.Errorf(t, err, msg.Error())
}

func TestIsSupportedType(t *testing.T) {
	wcTypes := GetSupportedWrkchainTypes()
	for _, wt := range wcTypes {
		isSupported := IsSupportedWrkchainType(wt)
		require.True(t, isSupported)
	}

	isSupported := IsSupportedWrkchainType("rubbishchain")
	require.False(t, isSupported)
}

func TestIsSupportedHash(t *testing.T) {
	wcTypes := GetSupportedWrkchainTypes()

	isSUpported, err := IsSupportedHash("neo", "MerkleRoot")
	require.Nil(t, err)
	require.True(t, isSUpported)

	isSUpported, err = IsSupportedHash("neo", "MerkleTrunk")
	require.Error(t, err)
	require.False(t, isSUpported)

	isSUpported, err = IsSupportedHash("garfunklechain", "MerkleTrunk")
	msg := fmt.Errorf("unknown wrkchain type rubbishchain, expected either %s", strings.Join(wcTypes, " or "))
	require.Errorf(t, err, msg.Error())
	require.False(t, isSUpported)
}

func TestGetDefaultHashMap(t *testing.T) {
	hashmap := GetDefaultHashMap("geth", "hash1")
	require.Equal(t, TxRoot, hashmap)
}

func TestGetSupportedWrkchainTypes(t *testing.T) {
	wcTypes := GetSupportedWrkchainTypes()

	isIn := contains(wcTypes, "geth")
	require.True(t, isIn)
	isIn = contains(wcTypes, "cosmos")
	require.True(t, isIn)
	isIn = contains(wcTypes, "tendermint")
	require.True(t, isIn)

	isIn = contains(wcTypes, "plum")
	require.False(t, isIn)
}

func TestGetLatestBlock(t *testing.T) {
	var wrkchainMeta = WrkChainMeta{
		LastBlock: "1",
		Type:      string(PseudoChainWrkchainType),
	}

	wc, err := NewWrkChain(wrkchainMeta, logger)
	require.Nil(t, err)

	wcBlock, err := wc.GetLatestBlock()
	require.Nil(t, err)
	require.Equal(t, true, len(wcBlock.BlockHash) > 0)

	wrkchainMeta.Type = string(EosWrkchainType)

	viper.Set(types.FlagWrkchainRpc, "https://api.testnet.eos.io")
	wc, err = NewWrkChain(wrkchainMeta, logger)
	require.Nil(t, err)

	wcBlock, err = wc.GetLatestBlock()
	require.Nil(t, err)
	require.Equal(t, true, len(wcBlock.BlockHash) > 0)
}

func TestGetWrkChainBlock(t *testing.T) {
	var wrkchainMeta = WrkChainMeta{
		LastBlock: "1",
		Type:      string(PseudoChainWrkchainType),
	}

	wc, err := NewWrkChain(wrkchainMeta, logger)
	require.Nil(t, err)

	wcBlock, err := wc.GetWrkChainBlock(1)
	require.Nil(t, err)
	require.Equal(t, true, len(wcBlock.BlockHash) > 0)

	wrkchainMeta.Type = string(EosWrkchainType)

	viper.Set(types.FlagWrkchainRpc, "https://api.testnet.eos.io")
	wc, err = NewWrkChain(wrkchainMeta, logger)
	require.Nil(t, err)

	wcBlock, err = wc.GetWrkChainBlock(1)
	require.Nil(t, err)
	require.Equal(t, "00000001c445c471d0b5563b0f73a113b14fcd3409df821c4ed42d32c0e3c5ca", wcBlock.BlockHash)
}
