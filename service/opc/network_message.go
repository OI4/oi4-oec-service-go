package opc

import (
	"fmt"
	"sync"
	"time"

	"github.com/OI4/oi4-oec-service-go/service/api"
)

var lastTimestamp int64
var counter int
var messageIDMutex = sync.Mutex{}

func GetMessageID(publisherID string) string {
	messageIDMutex.Lock()
	defer messageIDMutex.Unlock()

	currentTimestamp := time.Now().UnixMilli()

	if currentTimestamp <= lastTimestamp {
		counter++
	} else {
		counter = 0
		lastTimestamp = currentTimestamp
	}

	if counter == 0 {
		return fmt.Sprintf("%d-%s", currentTimestamp, publisherID)
	}
	return fmt.Sprintf("%d-%d-%s", currentTimestamp, counter, publisherID)
}

// CreateNetworkMessage quick and dirty
func CreateNetworkMessage(applicationOi4Identifier *api.Oi4Identifier, serviceType api.ServiceType, publication api.PublicationMessage) *api.NetworkMessage {
	content := publication.Content
	if content == nil || len(content) == 0 {
		return nil
	}

	source := publication.Source
	resourceType := publication.Resource
	assetOi4Identifier := publication.Source
	correlationId := publication.CorrelationId
	//payload := publication.Data

	datasetWriterId := GetDataSetWriterId(publication.Resource, source)

	currentTime := time.Now().UTC()

	messages := make([]*api.DataSetMessage, len(content))

	for i, message := range content {
		messages[i] = getMessageFromPayload(currentTime, datasetWriterId, applicationOi4Identifier, assetOi4Identifier, message)
	}

	networkMessage := &api.NetworkMessage{
		MessageId:      GetMessageID(applicationOi4Identifier.ToString()),
		MessageType:    api.UA_DATA,
		PublisherId:    fmt.Sprintf("%s/%s", serviceType, applicationOi4Identifier.ToString()),
		DataSetClassId: resourceType.ToDataSetClassId(),
		Messages:       messages,
		CorrelationId:  correlationId,
	}

	return networkMessage

}

func getMessageFromPayload(ts time.Time, datasetWriterId uint16, applicationOi4Identifier *api.Oi4Identifier, assetOi4Identifier *api.Oi4Identifier, content api.PublicationContent) *api.DataSetMessage {
	timestamp := ts.Format(time.RFC3339)

	message := &api.DataSetMessage{
		Timestamp:       &timestamp,
		DataSetWriterId: datasetWriterId,
		Status:          content.StatusCode,
		Payload:         content.Data,
	}

	if assetOi4Identifier != nil {
		message.Source = assetOi4Identifier.ToString()
	} else {
		message.Source = applicationOi4Identifier.ToString()
	}
	return message
}
