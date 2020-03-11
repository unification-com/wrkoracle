package oracle

import (
	"os"
	"time"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/unification-com/wrkoracle/client/mainchain"
	"github.com/unification-com/wrkoracle/client/wrkchains"
	"github.com/unification-com/wrkoracle/types"
)

// WrkOracle is an object to hold Oracle settings, and the WRKChain and Mainchain clients
type WrkOracle struct {
	frequency       uint64
	log             log.Logger
	wrkChain        *wrkchains.WrkChain
	mainchainClient *mainchain.MainchainClient
}

// NewWrkOracle returns an initialised WrkOracle object
func NewWrkOracle(cliCxt context.CLIContext, kb keys.Keybase, cdc *codec.Codec) (WrkOracle, error) {

	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	mc := mainchain.NewMainchainClient(cliCxt, kb, cdc, logger)
	err := mc.SetWrkchainMetaData()
	if err != nil {
		return WrkOracle{}, err
	}
	wrkChain, err := wrkchains.NewWrkChain(mc.GetWrkchainMeta(), logger)

	if err != nil {
		return WrkOracle{}, err
	}

	return WrkOracle{
		frequency:       viper.GetUint64(types.FlagFrequency),
		log:             logger.With("pkg", "oracle"),
		mainchainClient: mc,
		wrkChain:        wrkChain,
	}, nil
}

// Run runs the WRKOracle in automated mode
func (wo WrkOracle) Run() error {
	return wo.runOracle()
}

func (wo WrkOracle) runOracle() error {

	wo.log.Info("Start running WRKOracle")

	errors := make(chan error)

	for {
		go func() {
			timeNow := time.Now().Local()
			wo.log.Info("start poll", "time", timeNow)
			dueAt := time.Now().Local().Add(time.Duration(wo.frequency) * time.Second)
			err := wo.poll()
			if err != nil {
				errors <- err
				return
			}
			wo.log.Info("end poll. Next poll due:", "due", dueAt)
			wo.log.Info("-----------------------------------")
		}()
		select {
		case err := <-errors:
			wo.log.Error(err.Error())
			return err
		case <-time.After(time.Duration(wo.frequency) * time.Second):
		}
	}
}

func (wo WrkOracle) poll() error {

	wo.log.Info("polling WRKChain for latest block")
	header, err := wo.wrkChain.GetLatestBlock()

	if err != nil {
		wo.log.Error(err.Error())
		return err
	}

	wo.mainchainClient.SetRecordFees()

	wo.log.Info("recording latest WRKChain block")
	return wo.mainchainClient.BroadcastToMainchain(header)
}

// RecordSingleBlock is used to record a single WRKChain block header to Mainchain
func (wo WrkOracle) RecordSingleBlock(height uint64) error {
	return wo.recordBlock(height)
}

func (wo WrkOracle) recordBlock(height uint64) error {

	wo.log.Info("getting requested WRKChain block header and recording", "moniker", wo.mainchainClient.GetWrkchainMeta().Moniker, "height", height)

	header, err := wo.wrkChain.GetWrkChainBlock(height)

	if err != nil {
		wo.log.Error(err.Error())
		return err
	}

	wo.mainchainClient.SetRecordFees()

	return wo.mainchainClient.BroadcastToMainchain(header)
}
