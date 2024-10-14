package api

type ServiceParameterRequest struct {
	MethodsToCall []interface{} `json:"MessagesToCall"`
}
