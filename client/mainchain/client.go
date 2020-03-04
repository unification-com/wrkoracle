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

type MainchainClient struct {
	wrkchainID    uint64
	mainchainRest string
	wrkchainMeta  types.WrkChainMeta
	cliCtx        context.CLIContext
	kb            keys.Keybase
	cdc           *codec.Codec
	recFee        string
}

func NewMainchainClient(wrkchainID uint64, cliCtx context.CLIContext, kb keys.Keybase, cdc *codec.Codec) MainchainClient {
	mainchainRest := viper.GetString(types.FlagMainchainRest)
	return MainchainClient{
		wrkchainID:    wrkchainID,
		wrkchainMeta:  types.WrkChainMeta{},
		mainchainRest: mainchainRest,
		cliCtx:        cliCtx,
		kb:            kb,
		cdc:           cdc,
		recFee:        "1000000000nund",
	}
}

func (mc MainchainClient) BroadcastToMainchain(header types.WrkChainBlockHeader) error {

	fmt.Println("Generate msg")

	msg := wrkchain.NewMsgRecordWrkChainBlock(mc.wrkchainID, header.Height, header.BlockHash, header.ParentHash, header.Hash1, header.Hash2, header.Hash3, mc.cliCtx.GetFromAddress())
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

func (mc MainchainClient) GetWrkchainType() string {
	return mc.wrkchainMeta.Type
}

func (mc MainchainClient) GetWrkchainMeta() types.WrkChainMeta {
	return mc.wrkchainMeta
}

func (mc MainchainClient) GetRecordFees() string {
	return mc.recFee
}

func (mc *MainchainClient) SetWrkchainMetaData() error {

	if len(mc.wrkchainMeta.Type) == 0 {
		queryUrl := mc.mainchainRest + "/wrkchain/" + strconv.FormatUint(mc.wrkchainID, 10)

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

		mc.wrkchainMeta = wc.Result
		return nil
	}

	return nil
}

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
