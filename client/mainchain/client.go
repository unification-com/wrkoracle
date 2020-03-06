package mainchain

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/input"
	clikeys "github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/viper"
	"github.com/unification-com/mainchain/x/wrkchain"
	"github.com/unification-com/wrkoracle/types"
)

// MainchainClient is an object which holds data required to communicate with Mainchain
type MainchainClient struct {
	wrkchainId    uint64
	mainchainRest string
	wrkchainMeta  types.WrkChainMeta
	cliCtx        context.CLIContext
	kb            keys.Keybase
	cdc           *codec.Codec
	recFee        string
}

// NewMainchainClient returns an initialised MainchainClient object
func NewMainchainClient(wrkchainID uint64, cliCtx context.CLIContext, kb keys.Keybase, cdc *codec.Codec) MainchainClient {
	mainchainRest := viper.GetString(types.FlagMainchainRest)
	wrkchainType := viper.GetString(types.FlagWrkchainType)

	return MainchainClient{
		wrkchainId:    wrkchainID,
		wrkchainMeta:  types.WrkChainMeta{
			Type: wrkchainType,
		},
		mainchainRest: mainchainRest,
		cliCtx:        cliCtx,
		kb:            kb,
		cdc:           cdc,
		recFee:        "1000000000nund",
	}
}

// BroadcastToMainchain generates and broadcasts a TX containing a MsgRecordWrkChainBlock
// message to Mainchain
func (mc MainchainClient) BroadcastToMainchain(header types.WrkChainBlockHeader) error {

	fmt.Println("Generate msg")

	msg := wrkchain.NewMsgRecordWrkChainBlock(mc.wrkchainId, header.Height, header.BlockHash, header.ParentHash, header.Hash1, header.Hash2, header.Hash3, mc.cliCtx.GetFromAddress())
	err := msg.ValidateBasic()

	if err != nil {
		return err
	}

	fmt.Println("Broadcasting Tx and waiting for response...")

	res, err := mc.txBroadcaster(mc.cliCtx, []sdk.Msg{msg})

	if err != nil {
		return err
	}

	return mc.parseTsRes(res)
}

func (mc MainchainClient) txBroadcaster(cliCtx context.CLIContext, msgs []sdk.Msg) (sdk.TxResponse, error) {

	fmt.Println(fmt.Sprintf("WRKChain header hash recording fee: %s", mc.recFee))
	feeCoin, err := sdk.ParseCoin(mc.recFee)

	if err != nil {
		return sdk.TxResponse{}, err
	}

	fees := sdk.NewCoins(feeCoin)

	txBldr := auth.NewTxBuilder(utils.GetTxEncoder(mc.cdc), 0, 0, 0, 1.5, true, viper.GetString(flags.FlagChainID), "WRKOracle", fees, sdk.DecCoins{})
	txBldr, err = utils.PrepareTxBuilder(txBldr, cliCtx)

	if err != nil {
		return sdk.TxResponse{}, err
	}

	fromName := cliCtx.GetFromName()

	if txBldr.SimulateAndExecute() || cliCtx.Simulate {
		txBldr, err = utils.EnrichWithGas(txBldr, cliCtx, msgs)
		if err != nil {
			return sdk.TxResponse{}, err
		}

		gasEst := utils.GasEstimateResponse{GasEstimate: txBldr.Gas()}
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", gasEst.String())
	}

	if cliCtx.Simulate {
		return sdk.TxResponse{}, nil
	}

	if !cliCtx.SkipConfirm {
		stdSignMsg, err := txBldr.BuildSignMsg(msgs)
		if err != nil {
			return sdk.TxResponse{}, err
		}

		var jsonBz []byte
		if viper.GetBool(flags.FlagIndentResponse) {
			jsonBz, err = cliCtx.Codec.MarshalJSONIndent(stdSignMsg, "", "  ")
			if err != nil {
				panic(err)
			}
		} else {
			jsonBz = cliCtx.Codec.MustMarshalJSON(stdSignMsg)
		}

		_, _ = fmt.Fprintf(os.Stderr, "%s\n\n", jsonBz)

		buf := bufio.NewReader(os.Stdin)
		ok, err := input.GetConfirmation("confirm transaction before signing and broadcasting", buf)
		if err != nil || !ok {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", "cancelled transaction")
			return sdk.TxResponse{}, err
		}
	}

	// build and sign the transaction
	txBytes, err := txBldr.BuildAndSign(fromName, clikeys.DefaultKeyPass, msgs)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	// broadcast to a Tendermint node
	res, err := cliCtx.BroadcastTx(txBytes)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	return res, nil
}

func (mc MainchainClient) parseTsRes(res sdk.TxResponse) error {

	fmt.Println(fmt.Sprintf("Tx Hash: %s", res.TxHash))

	if len(res.Codespace) > 0 && res.Code > 0 {
		return fmt.Errorf("TX ERROR! Codespace: %s, Code: %d, Message: %s", res.Codespace, res.Code, res.RawLog)
	}

	fmt.Println(fmt.Sprintf("Success! Recorded in Mainchain Block #%d", res.Height))
	fmt.Println(fmt.Sprintf("Gas used: %d", res.GasUsed))

	return nil
}

// GetWrkchainType returns the configured WRKChain type, e.g. geth etc.
func (mc MainchainClient) GetWrkchainType() string {
	return mc.wrkchainMeta.Type
}

// GetWrkchainMeta returns the WRKChain's metadata object
func (mc MainchainClient) GetWrkchainMeta() types.WrkChainMeta {
	return mc.wrkchainMeta
}

// GetRecordFees returns the current required fees to submit WRKChain block header hashes to Mainchain
func (mc MainchainClient) GetRecordFees() string {
	return mc.recFee
}

// SetWrkchainMetaData queries Mainchain for the WRKChain metadata and stores it
// in the MainchainClient.wrkchainMeta object
func (mc *MainchainClient) SetWrkchainMetaData() error {

	if len(mc.wrkchainMeta.Moniker) == 0 {
		fmt.Println("Check WRKChain metadata")
		queryUrl := mc.mainchainRest + "/wrkchain/" + strconv.FormatUint(mc.wrkchainId, 10)

		resp, err := http.Get(queryUrl)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		var wc types.WrkChainMetaQueryResponse
		err = json.Unmarshal(body, &wc)
		if err != nil {
			return err
		}

		wrkchainType := wc.Result.Type

		if wrkchainType != mc.wrkchainMeta.Type {
			return fmt.Errorf("WRKChain Type mismatch: configured = %s, Mainchain = %s", mc.wrkchainMeta.Type, wrkchainType)
		}
		onChainId, err := strconv.Atoi(wc.Result.WRKChainId)
		if err != nil {
			return err
		}

		if uint64(onChainId) != mc.wrkchainId {
			return fmt.Errorf("WRKChain ID mismatch: configured = %s, Mainchain = %s", mc.wrkchainId, onChainId)
		}

		mc.wrkchainMeta = wc.Result
		return nil
	}

	return nil
}

// SetRecordFees queries Mainchain's wrkchain/params endpoint to update the current fees
// required to pay for submitting WRKChain block header hashes. The result is used internally
// by the MainchainClient.txBroadcaster function
func (mc *MainchainClient) SetRecordFees() {
	queryUrl := mc.mainchainRest + "/wrkchain/params"
	resp, err := http.Get(queryUrl)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var fees types.FeeParamsQueryResponse
	err = json.Unmarshal(body, &fees)
	if err == nil {
		fee := fees.Result.FeeRecord + fees.Result.Denom
		if fee != mc.recFee {
			fmt.Println(fmt.Sprintf("Fees for recording updated: %s", fee))
			mc.recFee = fee
		}
	}
}
