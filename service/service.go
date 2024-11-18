package service

import (
	"fmt"
	"github.com/OI4/oi4-oec-service-go/service/api"
)

func GetLocalizedText() string {
	request := api.CallMethodRequest{
		MethodId: "GetLocalizedText",
	}
	return fmt.Sprintf("MethodId: %s", request.MethodId)
}
