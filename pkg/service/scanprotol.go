package service

// TxPrepareData is the data part, those are the struct for ontology mobile scan protocol
// ignore the inner struct
type TxPrepareData struct {
	Params *Params `json:"params"`
	Action string  `json:"action"`
}
type Args struct {
	Value IntOrString `json:"value"`
	Name  string      `json:"name"`
}
type Functions struct {
	Operation string `json:"operation"`
	Args      []Args `json:"args"`
}
type InvokeConfig struct {
	GasPrice     uint64      `json:"gasPrice"`
	Payer        string      `json:"payer"`
	Functions    []Functions `json:"functions"`
	ContractHash string      `json:"contractHash"`
	GasLimit     uint64      `json:"gasLimit"`
}
type Params struct {
	InvokeConfig *InvokeConfig `json:"invokeConfig"`
}

// TxPrepare is the outer object
type TxPrepare struct {
	Desc      Desc   `json:"desc"`
	Signer    string `json:"signer"`
	Exp       int64  `json:"exp"`
	ID        string `json:"id"`
	Callback  string `json:"callback"`
	Signature string `json:"signature"`
	Data      string `json:"data"`
	Chain     string `json:"chain"`
	Ver       string `json:"ver"`
	Requester string `json:"requester"`
}

// Desc is description for this transaction
type Desc struct {
	Price  string `json:"price"`
	Detail string `json:"detail"`
	Type   string `json:"type"`
}
