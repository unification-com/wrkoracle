package oracle

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/spf13/viper"
	"github.com/unification-com/wrkoracle/client/mainchain"
	"github.com/unification-com/wrkoracle/client/wrkchains"
	"github.com/unification-com/wrkoracle/types"
)

type WrkOracle struct {
	wrkchainId  uint64
	wrkchainRpc string
	frequency   uint64
	cliCxt      context.CLIContext
	kb          keys.Keybase
	cdc         *codec.Codec
}

func NewWrkOracle(cliCxt context.CLIContext, kb keys.Keybase, cdc *codec.Codec) WrkOracle {

	wrkchainId := viper.GetUint64(types.FlagWrkChainId)
	frequency := viper.GetUint64(types.FlagFrequency)
	wrkchainRpc := viper.GetString(types.FlagWrkchainRpc)

	return WrkOracle{
		wrkchainId:  wrkchainId,
		frequency:   frequency,
		wrkchainRpc: wrkchainRpc,
		cliCxt:      cliCxt,
		kb:          kb,
		cdc:         cdc,
	}
}

func (wo WrkOracle) Run() error {
	return wo.runOracle()
}

func (wo WrkOracle) runOracle() error {

	fmt.Println("running")

	mc := mainchain.NewMainchainClient(wo.wrkchainId, wo.cliCxt, wo.kb, wo.cdc)
	err := mc.SetWrkchainMetaData()
	if err != nil {
		return err
	}

	errors := make(chan error)

	for {
		go func() {
			timeNow :=  time.Now().Local()
			fmt.Println(fmt.Sprintf("starting %s", timeNow))
			dueAt := time.Now().Local().Add(time.Duration(wo.frequency) * time.Second)
			err = wo.poll(&mc)
			if err != nil {
				errors <- err
				return
			}
			fmt.Println(fmt.Sprintf("Done. Next poll due at %s", dueAt))
			fmt.Println("-----------------------------------")
		}()
		select {
		case err := <-errors:
			return err
		case <-time.After(time.Duration(wo.frequency) * time.Second):
		}
	}
}

func (wo WrkOracle) poll(mc *mainchain.MainchainClient) error {

	fmt.Println("polling WRKChain for latest block")
	header, err := wrkchains.GetLatestBlock(mc.GetWrkchainMeta())

	if err != nil {
		return err
	}

	mc.SetRecordFees()

	fmt.Println("recording latest WRKChain block")
	return mc.BroadcastToMainchain(header)
}

func (wo WrkOracle) RecordSingleBlock(height uint64) error {
	return wo.recordBlock(height)
}

func (wo WrkOracle) recordBlock(height uint64) error {

	mc := mainchain.NewMainchainClient(wo.wrkchainId, wo.cliCxt, wo.kb, wo.cdc)
	err := mc.SetWrkchainMetaData()
	if err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("getting WRKChain '%s' block %d and recording", mc.GetWrkchainMeta().Moniker, height))

	header, err := wrkchains.GetWrkChainBlock(mc.GetWrkchainMeta(), height)

	if err != nil {
		return err
	}

	mc.SetRecordFees()

	return mc.BroadcastToMainchain(header)
}
