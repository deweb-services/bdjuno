package actions

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/forbole/bdjuno/v2/cmd/actions/types"
)

func handler(w http.ResponseWriter, r *http.Request) {

	// set the response header as JSON
	w.Header().Set("Content-Type", "application/json")

	// read request body
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	// parse the body as action payload
	var actionPayload types.ActionPayload
	err = json.Unmarshal(reqBody, &actionPayload)
	if err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	fmt.Println(actionPayload.Input.Arg1.Address)

	// Send the request params to the Action's generated handler function
	result, err := Account_balances(actionPayload.Input)

	// throw if an error happens
	if err != nil {
		errorObject := types.GraphQLError{
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
