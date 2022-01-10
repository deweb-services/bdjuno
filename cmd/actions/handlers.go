package actions

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	actionstypes "github.com/forbole/bdjuno/v2/cmd/actions/types"
	"github.com/forbole/bdjuno/v2/cmd/types"

	"github.com/forbole/bdjuno/v2/database"
	"github.com/forbole/bdjuno/v2/modules"
	"github.com/forbole/bdjuno/v2/types/config"
	"github.com/forbole/juno/v2/cmd/parse"
	junoconfig "github.com/forbole/juno/v2/types/config"
)

func accountBalancesHandler(w http.ResponseWriter, r *http.Request) {

	// set the response header as JSON
	w.Header().Set("Content-Type", "application/json")

	// read request body
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	// parse the body as action payload
	var actionPayload actionstypes.AccountBalancesPayload
	err = json.Unmarshal(reqBody, &actionPayload)
	if err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	parseCtx, sources, err := getCtxAndSources()
	if err != nil {
		http.Error(w, "error while getting parsing context & sources", http.StatusInternalServerError)
		return
	}

	// Send the request params to the Action's generated handler function
	result, err := AccountBalances(actionPayload.Input, sources, parseCtx)

	// throw if an error happens
	if err != nil {
		errorObject := actionstypes.GraphQLError{
			Message: err.Error(),
		}
		errorBody, _ := json.Marshal(errorObject)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorBody)
		return
	}

	// Write the response as JSON
	data, _ := json.Marshal(result)
	w.Write(data)
}

// func totalSupplyHandler(w http.ResponseWriter, r *http.Request) {

// 	// set the response header as JSON
// 	w.Header().Set("Content-Type", "application/json")

// 	// read request body
// 	reqBody, err := ioutil.ReadAll(r.Body)
// 	if err != nil {
// 		http.Error(w, "invalid payload", http.StatusBadRequest)
// 		return
// 	}

// 	// parse the body as action payload
// 	var actionPayload actionstypes.ActionPayload
// 	err = json.Unmarshal(reqBody, &actionPayload)
// 	if err != nil {
// 		http.Error(w, "invalid payload", http.StatusBadRequest)
// 		return
// 	}

// 	parseCfg := parse.NewConfig().
// 		WithDBBuilder(database.Builder).
// 		WithEncodingConfigBuilder(config.MakeEncodingConfig(types.GetBasicManagers())).
// 		WithRegistrar(modules.NewRegistrar(types.GetAddressesParser()))

// 	parseCtx, err := parse.GetParsingContext(parseCfg)
// 	if err != nil {
// 		http.Error(w, "error while getting parsing context", http.StatusInternalServerError)
// 		return
// 	}

// 	// Send the request params to the Action's generated handler function
// 	result, err := Account_balances(actionPayload.Input, parseCtx)

// 	// throw if an error happens
// 	if err != nil {
// 		errorObject := actionstypes.GraphQLError{
// 			Message: err.Error(),
// 		}
// 		errorBody, _ := json.Marshal(errorObject)
// 		w.WriteHeader(http.StatusBadRequest)
// 		w.Write(errorBody)
// 		return
// 	}

// 	// Write the response as JSON
// 	data, _ := json.Marshal(result)
// 	w.Write(data)
// }

func getCtxAndSources() (*parse.Context, *modules.Sources, error) {
	parseCfg := parse.NewConfig().
		WithDBBuilder(database.Builder).
		WithEncodingConfigBuilder(config.MakeEncodingConfig(types.GetBasicManagers())).
		WithRegistrar(modules.NewRegistrar(types.GetAddressesParser()))

	parseCtx, err := parse.GetParsingContext(parseCfg)
	if err != nil {
		return nil, nil, err
	}

	sources, err := modules.BuildSources(junoconfig.Cfg.Node, parseCtx.EncodingConfig)
	if err != nil {
		return nil, nil, err
	}

	return parseCtx, sources, nil
}
