package entity

type JsonResult struct {
	ErrorType int         `json:"errorType"`
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
}
