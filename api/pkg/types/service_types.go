package types

type ServiceType string

const (
	ServiceType_Registry     ServiceType = "Registry"
	ServiceType_OTConnector  ServiceType = "OTConnector"
	ServiceType_Utility      ServiceType = "Utility"
	ServiceType_Persistence  ServiceType = "Persistence"
	ServiceType_Aggregation  ServiceType = "Aggregation"
	ServiceType_OOCConnector ServiceType = "OOCConnector"
	ServiceType_ITConnector  ServiceType = "ITConnector"
)
