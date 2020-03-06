package config

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/unification-com/wrkoracle/types"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client/flags"
	toml "github.com/pelletier/go-toml"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagGet = "get"
)

var configDefaults = map[string]string{
	"chain-id":        "",
	"keyring-backend": "os",
	"output":          "json",
	"node":            "",
	"broadcast-mode":  "block",
	"wrkchain-id":     "",
	"frequency":       "60",
	"wrkchain-rpc":    "",
	"mainchain-rest":  "",
	"from":            "",
	"trust-node":      "false",
	"indent":          "true",
	"parent-hash":     "true",
	"hash1":           "",
	"hash2":           "",
	"hash3":           "",
	"wrkchain-type":   "",
}

// ConfigCmd returns a CLI command to interactively create an application CLI
// config file.
func ConfigCmd(defaultCLIHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config <key> [value]",
		Short: "Create or query an application CLI configuration file",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Create or query config values for WRKOracle

Output all config values:
$ %s config

Output a single value:
$ %s config mainchain-rest --get

Set a config value:
$ %s config parent-hash true
`,
				version.ClientName, version.ClientName, version.ClientName),
		),
		RunE: runConfigCmd,
		Args: cobra.RangeArgs(0, 2),
	}

	cmd.Flags().String(flags.FlagHome, defaultCLIHome,
		"set client's home directory for configuration")
	cmd.Flags().Bool(flagGet, false,
		"print configuration value or its default if unset")
	return cmd
}

func runConfigCmd(cmd *cobra.Command, args []string) error {
	cfgFile, err := ensureConfFile(viper.GetString(flags.FlagHome))
	if err != nil {
		return err
	}

	getAction := viper.GetBool(flagGet)
	if getAction && len(args) != 1 {
		return fmt.Errorf("wrong number of arguments")
	}

	// load configuration
	tree, err := loadConfigFile(cfgFile)
	if err != nil {
		return err
	}

	// print the config and exit
	if len(args) == 0 {
		return getAll(tree)
	}

	key := args[0]

	// get config value for a given key
	if getAction {
		return getConfValue(key, tree)
	}

	if len(args) != 2 {
		return fmt.Errorf("wrong number of arguments")
	}

	value := args[1]

	// set config value for a given key
	err = setConfValue(tree, key, value)

	if err != nil {
		return err
	}

	// save configuration to disk
	if err := saveConfigFile(cfgFile, tree); err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "configuration saved to %s\n", cfgFile)
	return nil
}

func setConfValue(tree *toml.Tree, key, value string) error {
	switch key {
	case "chain-id", "output", "node", "broadcast-mode", "keyring-backend",
		"wrkchain-id", "frequency", "wrkchain-rpc", "mainchain-rest", "from":
		tree.Set(key, value)

	case "wrkchain-type":
		if !types.IsSupportedWrkchainType(value) {
			supportedTypes := strings.Join(types.SupportedWrkchainTypes, ", ")
			return fmt.Errorf("unsupported WRKChain type: %s. supported types: %s", value, supportedTypes)
		}
		tree.Set(key, value)

	case "hash1", "hash2", "hash3":
		currentType := tree.Get("wrkchain-type").(string)
		if !types.IsSupportedWrkchainType(currentType) {
			supportedTypes := strings.Join(types.SupportedWrkchainTypes, ", ")
			return fmt.Errorf("unsupported WRKChain type detected in config: %s. supported types: %s", currentType, supportedTypes)
		}
		if !types.IsSupportedHash(currentType, value) {
			supportedHashes := strings.Join(types.SupportedHashMaps[currentType], ", ")
			return fmt.Errorf("unsupported hash map '%s' for wrkchain type %s. supported types: %s", value, currentType, supportedHashes)
		}
		tree.Set(key, value)

	case "trace", "trust-node", "indent", "parent-hash":
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}

		tree.Set(key, boolVal)

	default:
		return errUnknownConfigKey(key)
	}
	return nil
}

func getAll(tree *toml.Tree) error {
	s, err := tree.ToTomlString()
	if err != nil {
		return err
	}
	fmt.Print(s)
	return nil
}

func getConfValue(key string, tree *toml.Tree) error {
	switch key {
	case "trace", "trust-node", "indent", "parent-hash":
		fmt.Println(tree.GetDefault(key, false).(bool))

	default:
		if defaultValue, ok := configDefaults[key]; ok {
			fmt.Println(tree.GetDefault(key, defaultValue).(string))
			return nil
		}

		return errUnknownConfigKey(key)
	}

	return nil
}

// InitConfigCmd returns a CLI command to initialise a config file with default values
func InitConfigCmd(defaultCLIHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init [type]",
		Short: "Initialise a default config file",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Initialise WRKOracle with default config values. Config file is saved
 in %s/config/config.toml

If the config file already exixts, the command exits without creating defaults.

Example:
$ %s init geth
$ %s init tendermint
`, defaultCLIHome, version.ClientName, version.ClientName),
		),
		RunE: runInitConfigCmd,
		Args: cobra.ExactArgs(1),
	}

	cmd.Flags().String(flags.FlagHome, defaultCLIHome,
		"set client's home directory for configuration")
	return cmd
}

func runInitConfigCmd(cmd *cobra.Command, args []string) error {

	wrkchainType := args[0]

	if !types.IsSupportedWrkchainType(wrkchainType) {
		supportedTypes := strings.Join(types.SupportedWrkchainTypes, ", ")
		return fmt.Errorf("unsupported WRKChain type: %s. supported types: %s", wrkchainType, supportedTypes)
	}

	cfgFile, err := ensureConfFile(viper.GetString(flags.FlagHome))

	if err != nil {
		return err
	}

	// try loading configuration - bail out if it already exists
	tree, _ := loadConfigFile(cfgFile)
	if tree.Has("chain-id") {
		return fmt.Errorf("config file already exists at %s", cfgFile)
	}

	for k, v := range configDefaults {
		switch k {
		case "trace", "trust-node", "indent", "parent-hash":
			boolVal, err := strconv.ParseBool(v)
			if err != nil {
				return err
			}

			tree.Set(k, boolVal)

		default:
			tree.Set(k, v)
		}
	}

	tree.Set(types.FlagWrkchainType, wrkchainType)
	err = initChainType(wrkchainType, tree)
	if err != nil {
		return err
	}

	// save configuration to disk
	if err := saveConfigFile(cfgFile, tree); err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "configuration saved to %s\n", cfgFile)

	return nil
}

func initChainType(wrkchainType string, tree *toml.Tree) error {

	switch wrkchainType {
	case "geth":
		tree.Set("hash1", "ReceiptHash")
		tree.Set("hash2", "TxHash")
		tree.Set("hash3", "Root")
	case "tendermint", "cosmos":
		tree.Set("hash1", "DataHash")
		tree.Set("hash2", "AppHash")
		tree.Set("hash3", "ValidatorsHash")
	default:
		return fmt.Errorf("unsupported WRKChain type: %s", wrkchainType)
	}

	return nil

}

func ensureConfFile(rootDir string) (string, error) {
	cfgPath := path.Join(rootDir, "config")
	if err := os.MkdirAll(cfgPath, os.ModePerm); err != nil {
		return "", err
	}

	return path.Join(cfgPath, "config.toml"), nil
}

// CheckConfigFileExists checks if the config.toml file exists
func CheckConfigFileExists() error {
	cfgFile, err := ensureConfFile(viper.GetString(flags.FlagHome))
	if err != nil {
		return err
	}

	if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
		return fmt.Errorf("%s does not exist\n", cfgFile)
	}

	return nil
}

func loadConfigFile(cfgFile string) (*toml.Tree, error) {
	if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "%s does not exist\n", cfgFile)
		return toml.Load(``)
	}

	bz, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return nil, err
	}

	tree, err := toml.LoadBytes(bz)
	if err != nil {
		return nil, err
	}

	return tree, nil
}

func saveConfigFile(cfgFile string, tree io.WriterTo) error {
	fp, err := os.OpenFile(cfgFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer fp.Close()

	_, err = tree.WriteTo(fp)
	return err
}

func errUnknownConfigKey(key string) error {
	return fmt.Errorf("unknown configuration key: %q", key)
}
