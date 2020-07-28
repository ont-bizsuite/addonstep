package meta

// AddonConfigSteps contain the total steps to config this addon
type AddonConfigSteps struct {
	Steps []*Step `json:"steps"`
}

// Step define each step's detail
type Step struct {
	Index       int32                  `json:"index"`       // index of total steps, NOTE: start from 1, not 0
	Name        string                 `json:"name"`        // current step name
	Description string                 `json:"description"` // more verbose description
	Params      map[string]interface{} `json:"params"`      // request param
	Path        string                 `json:"path"`        // the relative HTTP path for current steps
	IsTx        bool                   `json:"tx"`          // is this step a transaction? addon server will use different way to process transaction step and none transaction step, see readme for detail description

	IsRollbackTx []bool   `json:"is_rollback_tx_lst"` // one step may need multiple rollback steps, the length of the array indicate the total rollback steps, and each value indicate wether the step is a transaction step
	RollbackPath []string `json:"rollback_path_lst"`  // rollback request HTTP path
}
