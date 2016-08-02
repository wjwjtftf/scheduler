package common

// job response wapper
type JobResponse struct {
	Success      bool        `json:"result"`
	Error        ErrorDTO    `json:"error"`
	Content      interface{} `json:"data"`
	ResponseTime int64       `json:"responseTime"`
	Status       string      `json:status`
}

type ErrorDTO struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
