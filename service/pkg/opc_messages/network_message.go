package opcmessages

import (
	"fmt"
	"time"

	v1 "github.com/OI4/oi4-oec-service-go/api/pkg/types"
)

var counter uint16 = 0

// CreateNetworkMessage quick and dirty
func CreateNetworkMessage(applicationOi4Identifier *v1.Oi4Identifier, serviceType v1.ServiceType, resourceType v1.ResourceType, assetOi4Identifier *v1.Oi4Identifier, datasetWriterId uint16, correlationId string, payload interface{}) *v1.NetworkMessage {
	currentTime := time.Now().UTC()

	message := &v1.DataSetMessage{
		Timestamp:       currentTime.Format(time.RFC3339),
		DataSetWriterId: datasetWriterId,
		Payload:         payload,
	}

	if assetOi4Identifier != nil {
		message.Source = assetOi4Identifier.ToString()
	} else {
		message.Source = applicationOi4Identifier.ToString()
	}

	networkMessage := &v1.NetworkMessage{
		MessageId:      fmt.Sprintf("%d%d-%s/%s", currentTime.Unix(), counter, serviceType, applicationOi4Identifier.ToString()),
		MessageType:    v1.UA_DATA,
		PublisherId:    fmt.Sprintf("%s/%s", serviceType, applicationOi4Identifier.ToString()),
		DataSetClassId: resourceType.ToDataSetClassId(),
		Messages:       []*v1.DataSetMessage{message},
		CorrelationId:  correlationId,
	}
	counter++

	return networkMessage

}
