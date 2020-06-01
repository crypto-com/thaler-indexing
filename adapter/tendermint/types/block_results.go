package types

type BlockResultsResp struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  struct {
		Height                string                     `json:"height"`
		TxsEvents             []RawBlockResultsTxResult  `json:"txs_results"`
		BeginBlockEvents      []RawBlockResultsEvent     `json:"begin_block_events"`
		EndBlockEvents        []RawBlockResultsEvent     `json:"end_block_events"`
		ValidatorUpdates      []RawBlockResultsValidator `json:"validator_updates"`
		ConsensusParamUpdates interface{}                `json:"consensus_param_updates"`
	} `json:"result"`
}

type RawBlockResultsTxResult struct {
	Code      int                    `json:"code"`
	Data      interface{}            `json:"data"`
	Log       string                 `json:"log"`
	Info      string                 `json:"info"`
	GasWanted string                 `json:"gasWanted"`
	GasUsed   string                 `json:"gasUsed"`
	Events    []RawBlockResultsEvent `json:"events"`
	Codespace string                 `json:"codespace"`
}

type RawBlockResultsEvent struct {
	Type       string `json:"type"`
	Attributes []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"attributes"`
}

type RawBlockResultsValidator struct {
	PubKey struct {
		Type string `json:"type"`
		Data string `json:"data"`
	} `json:"pub_key"`
	Power *string `json:"power"`
}

type BlockResults struct {
	Height           uint64
	TxsEvents        [][]BlockResultsEvent
	BeginBlockEvents []BlockResultsEvent
	ValidatorUpdates []BlockResultsValidator
}

type BlockResultsEvent struct {
	Type       string
	Attributes []BlockResultsEventAttribute
}
type BlockResultsEventAttribute struct {
	Key   string
	Value string
}

type BlockResultsValidator struct {
	PubKey BlockResultsValidatorPubKey
	Power  *string
}
type BlockResultsValidatorPubKey struct {
	Type    string
	PubKey  string
	Address string
}
