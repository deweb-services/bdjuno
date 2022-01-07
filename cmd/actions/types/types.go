package types

type ActionPayload struct {
	SessionVariables map[string]interface{} `json:"session_variables"`
	Input            Account_balancesArgs   `json:"input"`
}

type GraphQLError struct {
	Message string `json:"message"`
}

type Address struct {
	Address string
}

type Coins struct {
	Coins string
}

type Mutation struct {
	Account_balances *Coins
}

type Account_balancesArgs struct {
	Arg1 Address
}
