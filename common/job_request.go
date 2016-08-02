package common

// job 请求

const (
	INIT      = "INIT"
	EXECUTING = "EXECUTING"
	STOP      = "STOP"
	TEST      = "TEST"
	INVOKING  = "INVOKING"
	ERROR     = "ERROR"
)

type JobRequest struct {
	JobSnapshot int    `json:"jobSnapshot"`
	Params      string `json:"params"`
	Status      string `json:"status"`
}
