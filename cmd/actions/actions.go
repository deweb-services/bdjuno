package actions

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/forbole/bdjuno/v2/cmd/actions/types"
	"github.com/forbole/bdjuno/v2/database"
	"github.com/forbole/bdjuno/v2/modules"
	"github.com/forbole/bdjuno/v2/modules/bank"
	"github.com/forbole/juno/v2/modules/messages"

	"github.com/forbole/bdjuno/v2/types/config"
	junoconfig "github.com/forbole/juno/v2/types/config"

	"github.com/forbole/juno/v2/cmd/parse"
	"github.com/spf13/cobra"

	sifchainapp "github.com/Sifchain/sifnode/app"
	parsecmd "github.com/forbole/juno/v2/cmd/parse"
)

// NewActionsCmd returns the Cobra command allowing to activate hasura actions
func NewActionsCmd(parseCfg *parse.Config) *cobra.Command {
	return &cobra.Command{
		Use:     "hasura-actions",
		Short:   "Activate hasura actions",
		PreRunE: parse.ReadConfig(parseCfg),
		RunE: func(cmd *cobra.Command, args []string) error {

			// HTTP server for the handler
			mux := http.NewServeMux()
			mux.HandleFunc("/account_balances", handler)

			err := http.ListenAndServe(":3000", mux)
			log.Fatal(err)

			return nil
		},
	}
}

func Account_balances(args types.Account_balancesArgs) (response types.Coins, err error) {

	parseCfg := parsecmd.NewConfig().
		WithDBBuilder(database.Builder).
		WithEncodingConfigBuilder(config.MakeEncodingConfig(getBasicManagers())).
		WithRegistrar(modules.NewRegistrar(getAddressesParser()))

	parseCtx, err := parse.GetParsingContext(parseCfg)
	if err != nil {
		return types.Coins{}, err
	}
	sources, err := modules.BuildSources(junoconfig.Cfg.Node, parseCtx.EncodingConfig)
	if err != nil {
		return types.Coins{}, err
	}
	bankModule := bank.NewModule(nil, sources.BankSource, parseCtx.EncodingConfig.Marshaler, nil)

	// Get latest height
	height, err := parseCtx.Node.LatestHeight()
	if err != nil {
		return types.Coins{}, fmt.Errorf("error while getting chain latest block height: %s", err)
	}

	balances, err := bankModule.Keeper.GetBalances([]string{args.Arg1.Address}, height)
	if err != nil {
		return types.Coins{}, err
	}

	json, err := json.Marshal(balances)
	if err != nil {
		return types.Coins{}, err
	}

	response = types.Coins{
		Coins: string(json),
	}
	return response, nil
}

// getBasicManagers returns the various basic managers that are used to register the encoding to
// support custom messages.
// This should be edited by custom implementations if needed.
func getBasicManagers() []module.BasicManager {
	return []module.BasicManager{
		simapp.ModuleBasics,
		sifchainapp.ModuleBasics,
	}
}

// getAddressesParser returns the messages parser that should be used to get the users involved in
// a specific message.
// This should be edited by custom implementations if needed.
func getAddressesParser() messages.MessageAddressesParser {
	return messages.JoinMessageParsers(
		messages.CosmosMessageAddressesParser,
	)
}
