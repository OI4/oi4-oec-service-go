package api

import "fmt"

type ServiceType string

const (
	ServiceTypeRegistry     ServiceType = "Registry"
	ServiceTypeOTConnector  ServiceType = "OTConnector"
	ServiceTypeUtility      ServiceType = "Utility"
	ServiceTypePersistence  ServiceType = "Persistence"
	ServiceTypeAggregation  ServiceType = "Aggregation"
	ServiceTypeOOCConnector ServiceType = "OOCConnector"
	ServiceTypeITConnector  ServiceType = "ITConnector"
)

var serviceTypes = map[ServiceType]struct{}{
	ServiceTypeRegistry:     {},
	ServiceTypeOTConnector:  {},
	ServiceTypeUtility:      {},
	ServiceTypePersistence:  {},
	ServiceTypeAggregation:  {},
	ServiceTypeOOCConnector: {},
	ServiceTypeITConnector:  {},
}

func ParseServiceType(s string) (*ServiceType, error) {
	sType := ServiceType(s)
	_, ok := serviceTypes[sType]
	if !ok {
		return nil, fmt.Errorf(`cannot parse:[%s] as ServiceType`, s)
	}
	return &sType, nil
}
