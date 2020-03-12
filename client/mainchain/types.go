package mainchain

// FeeParams holds fee data from Mainchain params query for WRKChain fees
type FeeParams struct {
	Denom       string `json:"denom"`
	FeeRecord   string `json:"fee_record"`
	FeeRegister string `json:"fee_register"`
}

// FeeParamsQueryResponse is is a structure which holds Mainchain query data
// specifically for WRKChain module params queries
type FeeParamsQueryResponse struct {
	Height string    `json:"height"`
	Result FeeParams `json:"result"`
}
