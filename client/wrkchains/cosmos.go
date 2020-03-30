package wrkchains

import (
	"github.com/tendermint/tendermint/libs/log"
	"github.com/unification-com/wrkoracle/types"
)

func init() {
	wrkchainClientCreator := func(log log.Logger, lastHeight uint64) WrkChainClient {
		return NewTendermintClient(log, lastHeight, CosmosWrkchainType)
	}

	supportedHashMaps := []string{DataHash, AppHash, ValidatorsHash, LastResultsHash, LastCommitHash, ConsensusHash, NextValidatorsHash, EvidenceHash}

	defaultHashMap := make(map[string]string)
	defaultHashMap[types.FlagHash1] = DataHash
	defaultHashMap[types.FlagHash2] = AppHash
	defaultHashMap[types.FlagHash3] = ValidatorsHash
	registerWrkchainModule(CosmosWrkchainType, wrkchainClientCreator, supportedHashMaps, defaultHashMap, false)
}

// We can reuse the Tendermint struct
var _ WrkChainClient = (*Tendermint)(nil)
