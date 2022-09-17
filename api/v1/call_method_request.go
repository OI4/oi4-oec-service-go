package v1

type CallMethodRequest struct {
	MethodId       string        `json:"MethodId"`
	InputArguments []interface{} `json:"InputArguments"`
}
