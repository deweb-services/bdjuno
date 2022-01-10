package actions

import (
	"encoding/json"
	"fmt"

	actionstypes "github.com/forbole/bdjuno/v2/cmd/actions/types"
	"github.com/forbole/bdjuno/v2/modules"
	"github.com/forbole/bdjuno/v2/modules/bank"
	"github.com/forbole/juno/v2/cmd/parse"
)

func AccountBalances(args actionstypes.Account_balancesArgs, sources *modules.Sources, parseCtx *parse.Context) (response actionstypes.Coins, err error) {

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
