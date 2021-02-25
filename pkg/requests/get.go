package requests

type Get struct {
	BaseRequest
	Key string `json:"key"`
}
