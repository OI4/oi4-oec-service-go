package opcmessages

import (
	"fmt"
	"github.com/OI4/oi4-oec-service-go/api/pkg/types"
	"sync"
)

type dataSetWriterId struct {
	id        uint16
	writerIds map[string]uint16
	mutex     sync.RWMutex
}

var dswId = dataSetWriterId{
	id:        9,
	writerIds: make(map[string]uint16),
	mutex:     sync.RWMutex{},
}

func GetDataSetWriterId(resource types.ResourceType, source types.Oi4Identifier) uint16 {
	key := getDataSetWriterIdKey(resource, source)
	dswId.mutex.RLock()
	defer dswId.mutex.RUnlock()
	if _, ok := dswId.writerIds[key]; !ok {
		dswId.id += 1
		dswId.writerIds[key] = dswId.id
	}
	return dswId.writerIds[key]
}

func getDataSetWriterIdKey(resource types.ResourceType, source types.Oi4Identifier) string {
	sub := source.ToString()
	if resource == types.ResourcePublicationList || resource == types.ResourceSubscriptionList {
		sub = "NA"
	}
	return fmt.Sprintf("%s_|_%s", resource, sub)
}
