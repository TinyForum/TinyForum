package common

type ResponseMessage struct {
	Message string `json:"message"`
	Pass    bool   `json:"status"`
}
type ResponseBool struct {
	Success bool `json:"success"`
}
