package main

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	clikeys "github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/unification-com/mainchain/app"
	undtypes "github.com/unification-com/mainchain/types"
	"github.com/unification-com/wrkoracle/config"
	"github.com/unification-com/wrkoracle/oracle"
	"github.com/unification-com/wrkoracle/types"
)

var (
	defaultHome = os.ExpandEnv("$HOME/.und_wrkoracle")
)

func main() {
	cobra.EnableCommandSorting = false

	cdc := app.MakeCodec()

	// Read in the configuration file for the sdk
	sdkConfig := sdk.GetConfig()
	sdkConfig.SetBech32PrefixForAccount(undtypes.Bech32PrefixAccAddr, undtypes.Bech32PrefixAccPub)
	sdkConfig.SetBech32PrefixForValidator(undtypes.Bech32PrefixValAddr, undtypes.Bech32PrefixValPub)
	sdkConfig.SetBech32PrefixForConsensusNode(undtypes.Bech32PrefixConsAddr, undtypes.Bech32PrefixConsPub)
	sdkConfig.SetCoinType(undtypes.CoinType)
	sdkConfig.SetFullFundraiserPath(undtypes.HdWalletPath)
	sdkConfig.Seal()

	rootCmd := &cobra.Command{
		Use:   "wrkoracle",
		Short: "WRKOracle CLI tool for submitting WRKChain hashes to Mainchain",
	}

	rootCmd.PersistentFlags().String(flags.FlagChainID, "", "Chain ID of UND Mainchain")
	rootCmd.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
		return initConfig(rootCmd)
	}

	rootCmd.AddCommand(
		config.ConfigCmd(defaultHome),
		config.InitConfigCmd(defaultHome),
		version.Cmd,
	)
	rootCmd.AddCommand(
		flags.PostCommands(
			clikeys.Commands(),
			RunCmd(cdc),
			RecordSingleCmd(cdc),
		)...,
	)

	executor := cli.PrepareMainCmd(rootCmd, "UND", defaultHome)
	err := executor.Execute()
	if err != nil {
		fmt.Printf("Failed executing CLI command: %s, exiting...\n", err)
		os.Exit(1)
	}
}

// RunCmd is the CLI command for running the WRKOracle in automatic mode.
// It will run the oracle and poll the WRKChain according to the configured
// frequency.
func RunCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "run wrkoracle to record a WRKChain's block hashes",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Run WRKOracle to record a new WRKChain block's hash(es)`),
		),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {

			err := config.CheckConfigFileExists()
			if err != nil {
				return fmt.Errorf("wrkoracle not yet initialised. Run `wrkoracle init [wrkchain_type]`")
			}

			from := viper.GetString(flags.FlagFrom)
			wrkchainType := viper.GetString(types.FlagWrkchainType)

			// set --yes flag true, otherwise it cannot automate broadcasting Txs
			viper.Set(flags.FlagSkipConfirmation, true)

			if viper.GetUint(types.FlagWrkChainId) == 0 {
				return fmt.Errorf("missing WRKChain ID: set %s in %s/config/config.toml or pass with --%s flag", types.FlagWrkChainId, defaultHome, types.FlagWrkChainId)
			}
			if viper.GetUint(types.FlagFrequency) == 0 {
				return fmt.Errorf("frequency must be > 0: set %s in %s/config/config.toml or pass with --%s flag", types.FlagFrequency, defaultHome, types.FlagFrequency)
			}
			if len(viper.GetString(types.FlagWrkchainRpc)) <= 0 {
				return fmt.Errorf("missing WRKChain RPC URL: set %s in %s/config/config.toml or pass with --%s flag", types.FlagWrkchainRpc, defaultHome, types.FlagWrkchainRpc)
			}
			if len(viper.GetString(types.FlagMainchainRest)) <= 0 {
				return fmt.Errorf("missing Mainchain REST URL: set %s in %s/config/config.toml or pass with --%s flag", types.FlagMainchainRest, defaultHome, types.FlagMainchainRest)
			}
			if len(from) <= 0 {
				return fmt.Errorf("missing sender: set %s in %s/config/config.toml or pass with --%s flag", flags.FlagFrom, defaultHome, flags.FlagFrom)
			}
			if len(wrkchainType) <= 0 {
				return fmt.Errorf("missing WRKChain Type: set %s in %s/config/config.toml or pass with --%s flag", types.FlagWrkchainType, defaultHome, types.FlagWrkchainType)
			}

			if !types.IsSupportedWrkchainType(wrkchainType) {
				supportedTypes := strings.Join(types.SupportedWrkchainTypes, ", ")
				return fmt.Errorf("unsupported WRKChain type: %s. supported types: %s", wrkchainType, supportedTypes)
			}

			kb, err := keys.NewKeyring(sdk.KeyringServiceName(), viper.GetString(flags.FlagKeyringBackend), viper.GetString(flags.FlagHome), cmd.InOrStdin())

			if err != nil {
				return err
			}

			cliCtx := context.NewCLIContextWithInputAndFrom(cmd.InOrStdin(), from).WithCodec(cdc)

			wrkOracle, err := oracle.NewWrkOracle(cliCtx, kb, cdc)
			if err != nil {
				return err
			}
			return wrkOracle.Run()
		},
	}
	cmd.Flags().String(types.FlagWrkchainType, "", "WRKChain Type: geth, tendermint etc.")
	cmd.Flags().String(types.FlagWrkchainRpc, "", "WRKChain's RPC URL")
	cmd.Flags().Uint64(types.FlagFrequency, 0, "Frequency to submit WRKChain hashes in seconds")
	cmd.Flags().Uint64(types.FlagWrkChainId, 0, "WRKChain ID")
	cmd.Flags().String(types.FlagMainchainRest, "", "Mainchain REST URL")
	cmd.Flags().Bool(types.FlagParentHash, false, "submit parent hash")
	cmd.Flags().Bool(types.FlagHash1, false, "submit hash1")
	cmd.Flags().Bool(types.FlagHash2, false, "submit hash2")
	cmd.Flags().Bool(types.FlagHash3, false, "submit hash3")
	return cmd
}

// RecordSingleCmd can be used to record a single block's hashes for the given height
func RecordSingleCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "record [height]",
		Short: "use wrkoracle to record a WRKChain's single block hashes at the specified height",
		Long: strings.TrimSpace(
			fmt.Sprintf(`use wrkoracle to record a WRKChain's single block hashes at the specified height`),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			err := config.CheckConfigFileExists()
			if err != nil {
				return fmt.Errorf("wrkoracle not yet initialised. Run `wrkoracle init [wrkchain_type]`")
			}

			from := viper.GetString(flags.FlagFrom)
			wrkchainType := viper.GetString(types.FlagWrkchainType)

			// set --yes flag true, otherwise it cannot automate
			viper.Set(flags.FlagSkipConfirmation, true)

			if viper.GetUint(types.FlagWrkChainId) == 0 {
				return fmt.Errorf("missing WRKChain ID: set %s in %s/config/config.toml or pass with --%s flag", types.FlagWrkChainId, defaultHome, types.FlagWrkChainId)
			}
			if len(viper.GetString(types.FlagWrkchainRpc)) <= 0 {
				return fmt.Errorf("missing WRKChain RPC URL: set %s in %s/config/config.toml or pass with --%s flag", types.FlagWrkchainRpc, defaultHome, types.FlagWrkchainRpc)
			}
			if len(viper.GetString(types.FlagMainchainRest)) <= 0 {
				return fmt.Errorf("missing Mainchain REST URL: set %s in %s/config/config.toml or pass with --%s flag", types.FlagMainchainRest, defaultHome, types.FlagMainchainRest)
			}
			if len(from) <= 0 {
				return fmt.Errorf("missing sender: set %s in %s/config/config.toml or pass with --%s flag", flags.FlagFrom, defaultHome, flags.FlagFrom)
			}
			if len(wrkchainType) <= 0 {
				return fmt.Errorf("missing WRKChain Type: set %s in %s/config/config.toml or pass with --%s flag", types.FlagWrkchainType, defaultHome, types.FlagWrkchainType)
			}

			if !types.IsSupportedWrkchainType(wrkchainType) {
				supportedTypes := strings.Join(types.SupportedWrkchainTypes, ", ")
				return fmt.Errorf("unsupported WRKChain type: %s. supported types: %s", wrkchainType, supportedTypes)
			}

			kb, err := keys.NewKeyring(sdk.KeyringServiceName(), viper.GetString(flags.FlagKeyringBackend), viper.GetString(flags.FlagHome), cmd.InOrStdin())

			if err != nil {
				return err
			}

			height, err := strconv.Atoi(args[0])

			if err != nil {
				return err
			}

			cliCtx := context.NewCLIContextWithInputAndFrom(cmd.InOrStdin(), from).WithCodec(cdc)

			wrkOracle, err := oracle.NewWrkOracle(cliCtx, kb, cdc)
			if err != nil {
				return err
			}
			return wrkOracle.RecordSingleBlock(uint64(height))
		},
	}
	cmd.Flags().String(types.FlagWrkchainType, "", "WRKChain Type: geth, tendermint etc.")
	cmd.Flags().String(types.FlagWrkchainRpc, "", "WRKChain's RPC URL")
	cmd.Flags().Uint64(types.FlagWrkChainId, 0, "WRKChain ID")
	cmd.Flags().String(types.FlagMainchainRest, "", "Mainchain REST URL")
	cmd.Flags().Bool(types.FlagParentHash, false, "submit parent hash")
	cmd.Flags().Bool(types.FlagHash1, false, "submit hash1")
	cmd.Flags().Bool(types.FlagHash2, false, "submit hash2")
	cmd.Flags().Bool(types.FlagHash3, false, "submit hash3")
	return cmd
}

func initConfig(cmd *cobra.Command) error {
	home, err := cmd.PersistentFlags().GetString(cli.HomeFlag)
	if err != nil {
		return err
	}

	cfgFile := path.Join(home, "config", "config.toml")
	if _, err := os.Stat(cfgFile); err == nil {
		viper.SetConfigFile(cfgFile)

		if err := viper.ReadInConfig(); err != nil {
			return err
		}
	}
	if err := viper.BindPFlag(flags.FlagChainID, cmd.PersistentFlags().Lookup(flags.FlagChainID)); err != nil {
		return err
	}
	if err := viper.BindPFlag(cli.EncodingFlag, cmd.PersistentFlags().Lookup(cli.EncodingFlag)); err != nil {
		return err
	}
	return viper.BindPFlag(cli.OutputFlag, cmd.PersistentFlags().Lookup(cli.OutputFlag))
}
