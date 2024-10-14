package opc

import (
	"fmt"
	"time"

	"github.com/OI4/oi4-oec-service-go/service/api"
)

var counter uint16 = 0

// CreateNetworkMessage quick and dirty
func CreateNetworkMessage(applicationOi4Identifier *api.Oi4Identifier, serviceType api.ServiceType, resourceType api.ResourceType, assetOi4Identifier *api.Oi4Identifier, datasetWriterId uint16, correlationId string, payload interface{}) *api.NetworkMessage {
	currentTime := time.Now().UTC()

	message := &api.DataSetMessage{
		Timestamp:       currentTime.Format(time.RFC3339),
		DataSetWriterId: datasetWriterId,
		Payload:         payload,
	}

	if assetOi4Identifier != nil {
		message.Source = assetOi4Identifier.ToString()
	} else {
		message.Source = applicationOi4Identifier.ToString()
	}

	networkMessage := &api.NetworkMessage{
		MessageId:      fmt.Sprintf("%d%d-%s/%s", currentTime.Unix(), counter, serviceType, applicationOi4Identifier.ToString()),
		MessageType:    api.UA_DATA,
		PublisherId:    fmt.Sprintf("%s/%s", serviceType, applicationOi4Identifier.ToString()),
		DataSetClassId: resourceType.ToDataSetClassId(),
		Messages:       []*api.DataSetMessage{message},
		CorrelationId:  correlationId,
	}
	counter++

	return networkMessage

}
