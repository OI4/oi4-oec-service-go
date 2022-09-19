package types

type ServiceParameterRequest struct {
	MethodsToCall []interface{} `json:"MessagesToCall"`
}
