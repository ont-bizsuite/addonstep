package meta

import (
	"github.com/gin-gonic/gin"
	"github.com/ont-bizsuite/addonstep/pkg/path"
	"github.com/ont-bizsuite/addonstep/pkg/service"
)

// AddonConfigSteps contain the total steps to config this addon
type AddonConfigSteps struct {
	Steps []*Step `json:"steps"`
}

// Step define each step's detail
type Step struct {
	Index       int32       `json:"index"`       // index of total steps, NOTE: start from 1, not 0
	Name        string      `json:"name"`        // current step name
	Description string      `json:"description"` // more verbose description
	Params      interface{} `json:"params"`      // request param
	Path        string      `json:"path"`        // the relative HTTP path for current steps
	IsTx        bool        `json:"tx"`          // is this step a transaction? addon server will use different way to process transaction step and none transaction step, see readme for detail description

	IsRollbackTx []bool   `json:"is_rollback_tx_lst,omitempty"` // one step may need multiple rollback steps, the length of the array indicate the total rollback steps, and each value indicate wether the step is a transaction step
	RollbackPath []string `json:"rollback_path_lst,omitempty"`  // rollback request HTTP path

	AsyncPath string `json:"async_path"`
}

const (
	// StatusInit ...
	StatusInit = 1
	// StatusDone ...
	StatusDone = 2
	// StatusNE ...
	StatusNE = 3
)

type (
	// AsyncPathInput for async input
	AsyncPathInput struct {
		TenantID string `json:"tenant_id"`
		StepIdx  int    `json:"step_idx"`
	}

	// BaseResp ...
	BaseResp struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}

	// AsyncPathOutput for async output
	AsyncPathOutput struct {
		BaseResp
		Result interface{} `json:"result"`
		State  int         `json:"state"`
	}
)

var (
	steps = &AddonConfigSteps{}
)

var (
	StepPay = &Step{
		Name:        "Operating fee",
		Description: "Operating fee",
		Path:        path.PayPath,
		IsTx:        true,
		Params:      service.PaySample,
	}
)

// RegistPath is only for demo, call your own's addon register path!
func RegistPath(r *gin.Engine, sc *service.ServiceConfig) error {
	// meta
	r.GET(path.MetaSteps, AddonSteps)

	// service
	r.POST(path.PayPath, sc.TransferMoney)
	r.POST(path.PayCallbackPath, sc.FeeBack)

	return nil
}

func AddonSteps(c *gin.Context) {
	c.JSON(200, steps)
}

func AppendStep(s *Step) {
	// ensure steps
	defer func() {
		for i := range steps.Steps {
			steps.Steps[i].Index = int32(i) + 1
		}
	}()

	steps.Steps = append(steps.Steps, s)
}
