package v1

type ServiceParameterRequest struct {
	MethodsToCall []interface{} `json:"MessagesToCall"`
}
