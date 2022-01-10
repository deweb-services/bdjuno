package actions

import (
	"encoding/json"
	"fmt"

	actionstypes "github.com/forbole/bdjuno/v2/cmd/actions/types"
	"github.com/forbole/bdjuno/v2/modules"
	"github.com/forbole/bdjuno/v2/modules/bank"
	"github.com/forbole/juno/v2/cmd/parse"
	junoconfig "github.com/forbole/juno/v2/types/config"
)

func Account_balances(args actionstypes.Account_balancesArgs, parseCtx *parse.Context) (response actionstypes.Coins, err error) {

	sources, err := modules.BuildSources(junoconfig.Cfg.Node, parseCtx.EncodingConfig)
	if err != nil {
		return actionstypes.Coins{}, err
	}
	bankModule := bank.NewModule(nil, sources.BankSource, parseCtx.EncodingConfig.Marshaler, nil)

	// Get latest height
	height, err := parseCtx.Node.LatestHeight()
	if err != nil {
		return actionstypes.Coins{}, fmt.Errorf("error while getting chain latest block height: %s", err)
	}

	balances, err := bankModule.Keeper.GetBalances([]string{args.Arg1.Address}, height)
	if err != nil {
		return actionstypes.Coins{}, err
	}

	json, err := json.Marshal(balances)
	if err != nil {
		return actionstypes.Coins{}, err
	}

	response = actionstypes.Coins{
		Coins: string(json),
	}
	return response, nil
}
